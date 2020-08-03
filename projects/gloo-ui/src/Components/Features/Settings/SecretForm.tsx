import { ReactComponent as GreenPlus } from 'assets/small-green-plus.svg';
import {
  SoloFormInput,
  SoloFormTypeahead,
  TableFormWrapper
} from 'Components/Common/Form/SoloFormField';
import { Formik } from 'formik';
import {
  AwsSecret,
  AzureSecret,
  Secret,
  TlsSecret
} from 'proto/gloo/projects/gloo/api/v1/secret_pb';
import { ResourceRef } from 'proto/solo-kit/api/v1/ref_pb';
import React from 'react';
import { configAPI } from 'store/config/api';
import useSWR from 'swr';

// TODO: modify for use outside a table
// TODO: set one source of truth for column names/order

export interface SecretValuesType {
  secretResourceRef: ResourceRef.AsObject;
  awsSecret: AwsSecret.AsObject;
  azureSecret: AzureSecret.AsObject;
  tlsSecret: TlsSecret.AsObject;
  oAuthSecret: { clientSecret: string };
}

interface Props {
  secretKind: Secret.KindCase;
  onCreateSecret: (
    values: SecretValuesType,
    secretKind: Secret.KindCase
  ) => void;
}

export const SecretForm: React.FC<Props> = props => {
  const { secretKind, onCreateSecret } = props;

  const { data: namespacesList, error: listNamespacesError } = useSWR(
    'listNamespaces',
    configAPI.listNamespaces
  );
  const { data: podNamespace, error: podNamespaceError } = useSWR(
    'getPodNamespace',
    configAPI.getPodNamespace
  );
  if (!podNamespace || !namespacesList) {
    return <div>Loading...</div>;
  }

  const initialValues: SecretValuesType = {
    secretResourceRef: {
      name: '',
      namespace: podNamespace
    },
    awsSecret: {
      accessKey: '',
      secretKey: '',
      sessionToken: ''
    },
    azureSecret: {
      apiKeysMap: []
    },
    tlsSecret: {
      certChain: '',
      privateKey: '',
      rootCa: ''
    },
    oAuthSecret: {
      clientSecret: ''
    }
  };

  const createSecret = (values: SecretValuesType) => {
    onCreateSecret(values, secretKind);
  };

  return (
    <>
      <Formik<SecretValuesType>
        initialValues={initialValues}
        onSubmit={createSecret}>
        {({ handleSubmit }) => (
          <>
            <TableFormWrapper>
              <SoloFormInput
                testId={`${secretKind}-secret-name`}
                name='secretResourceRef.name'
                placeholder='Name'
              />
              <SoloFormTypeahead
                testId={`${secretKind}-secret-namespace`}
                name='secretResourceRef.namespace'
                placeholder='Namespace'
                defaultValue={podNamespace}
                presetOptions={namespacesList.map(ns => {
                  return { value: ns };
                })}
              />
            </TableFormWrapper>
            {secretKind === Secret.KindCase.AWS && <AwsSecretFields />}
            {secretKind === Secret.KindCase.AZURE && <AzureSecretFields />}
            {secretKind === Secret.KindCase.TLS && <TlsSecretFields />}
            {secretKind === Secret.KindCase.OAUTH && <OAuthSecretFields />}
            <TableFormWrapper>
              <span className='text-green-400 cursor-pointer hover:text-green-300'>
                <GreenPlus
                  className='fill-current'
                  data-testid={`${secretKind}-secret-green-plus`}
                  style={{ cursor: 'pointer' }}
                  onClick={() => handleSubmit()}
                />
              </span>
            </TableFormWrapper>
          </>
        )}
      </Formik>
    </>
  );
};

const AwsSecretFields: React.FC = () => {
  return (
    <TableFormWrapper>
      <SoloFormInput
        testId='aws-secret-accessKey'
        name='awsSecret.accessKey'
        placeholder='Access Key'
      />
      <SoloFormInput
        testId='aws-secret-secretKey'
        name='awsSecret.secretKey'
        placeholder='Secret Key'
      />
    </TableFormWrapper>
  );
};

const TlsSecretFields: React.FC = () => {
  return (
    <TableFormWrapper>
      <SoloFormInput
        testId='tls-secret-'
        name='tlsSecret.certChain'
        placeholder='Cert Chain'
      />
      <SoloFormInput
        testId='tls-secret-'
        name='tlsSecret.privateKey'
        placeholder='Private Key'
      />
      <SoloFormInput
        testId='tls-secret-'
        name='tlsSecret.rootCa'
        placeholder='Root Ca'
      />
    </TableFormWrapper>
  );
};

const AzureSecretFields: React.FC = () => {
  return (
    <TableFormWrapper>
      <SoloFormInput
        testId='azure-secret-apiKeysMap'
        name='azureSecret.apiKeysMap'
        placeholder='Api Keys'
      />
    </TableFormWrapper>
  );
};

const OAuthSecretFields: React.FC = () => {
  return (
    <TableFormWrapper>
      <SoloFormInput
        testId='oauth-secret-clientSecret'
        name='oAuthSecret.clientSecret'
        placeholder='Client Secret'
      />
    </TableFormWrapper>
  );
};
