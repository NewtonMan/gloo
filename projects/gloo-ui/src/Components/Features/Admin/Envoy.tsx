import styled from '@emotion/styled';
import { ReactComponent as EnvoyLogo } from 'assets/envoy-logo.svg';
import { ConfigDisplayer } from 'Components/Common/DisplayOnly/ConfigDisplayer';
import { TallyContainer } from 'Components/Common/DisplayOnly/TallyInformationDisplay';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';
import { SectionCard } from 'Components/Common/SectionCard';
import { Status } from 'proto/solo-projects/projects/grpcserver/api/v1/types_pb';
import { Status as SoloKitStatus } from 'proto/solo-kit/api/v1/status_pb';
import * as React from 'react';
import { envoyAPI } from 'store/envoy/api';
import { colors, healthConstants, soloConstants } from 'Styles';
import useSWR from 'swr';

const InsideHeader = styled.div`
  display: flex;
  justify-content: space-between;
  font-size: 18px;
  line-height: 22px;
  margin-bottom: 18px;
  color: ${colors.novemberGrey};
`;

const EnvoyLogoFullSize = styled(EnvoyLogo)`
  width: 35px !important;
  height: 35px;
  max-height: none !important;
`;

type ExpandableSectionProps = { isExpanded: boolean };
const ExpandableSection = styled.div`
  max-height: ${(props: ExpandableSectionProps) =>
    props.isExpanded ? '1000px' : '0px'};
  overflow: ${(props: { isExpanded: boolean }) =>
    props.isExpanded ? 'auto' : 'hidden'};
  transition: max-height ${soloConstants.transitionTime};
  color: ${colors.septemberGrey};
`;

const Link = styled.div`
  cursor: pointer;
  color: ${colors.seaBlue};
  font-size: 14px;
`;

interface Props {}
export const getHealth = (
  code: number
): SoloKitStatus.StateMap[keyof SoloKitStatus.StateMap] => {
  switch (code) {
    case Status.Code.ERROR:
      return healthConstants.Error.value;
    case Status.Code.OK:
      return healthConstants.Good.value;
    default:
      return healthConstants.Pending.value;
  }
};

export const Envoy = (props: Props) => {
  const { data: envoysList, error } = useSWR(
    'listEnvoys',
    envoyAPI.getEnvoyList
  );

  const [envoysOpen, setEnvoysOpen] = React.useState<boolean[]>([]);

  React.useEffect(() => {
    if (!!envoysList && !!envoysList.length) {
      setEnvoysOpen(envoysList.map(e => false));
    }
  }, [envoysList?.length]);

  if (!envoysList?.length) {
    return <div>No Envoy cofiguration found.</div>;
  }

  const toggleExpansion = (indexToggled: number) => {
    setEnvoysOpen(
      envoysOpen.map((isOpen, ind) => {
        if (ind !== indexToggled) {
          return false;
        }

        return !isOpen;
      })
    );
  };

  return (
    <>
      {envoysList.map((envoy, ind) => {
        const hasConfigDump = !!envoy.raw && envoy.raw.content.length > 0;
        return (
          <SectionCard
            key={envoy.name + ind}
            cardName={envoy.name}
            logoIcon={<EnvoyLogoFullSize />}
            headerSecondaryInformation={[]}
            health={getHealth(envoy!.status!.code)}
            healthMessage={'Envoy Status'}>
            {envoy!.status!.message !== '' && (
              <TallyContainer color='orange'>
                {envoy!.status!.message!}
              </TallyContainer>
            )}
            <InsideHeader>
              <div>Code Log (Read Only)</div>{' '}
              {hasConfigDump ? (
                <FileDownloadLink
                  fileName={envoy.raw!.fileName}
                  fileContent={envoy.raw!.content}
                />
              ) : (
                <div>---</div>
              )}
            </InsideHeader>
            {hasConfigDump ? (
              <>
                <ExpandableSection isExpanded={envoysOpen[ind]}>
                  {' '}
                  <ConfigDisplayer content={envoy.raw!.content} isJson />
                </ExpandableSection>
                <Link onClick={() => toggleExpansion(ind)}>
                  {envoysOpen[ind] ? 'Hide' : 'View'} Envoy Config
                </Link>
              </>
            ) : (
              <div>
                <i>Install Gloo with </i>
                <code>gloo.gatewayProxies.gatewayProxy.readConfig</code>{' '}
                <i>enabled to view Envoy config.</i>
              </div>
            )}
          </SectionCard>
        );
      })}
    </>
  );
};
