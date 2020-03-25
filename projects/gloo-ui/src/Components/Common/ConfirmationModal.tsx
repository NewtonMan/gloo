import * as React from 'react';

import { Modal } from 'antd';
import { ModalProps } from 'antd/lib/modal';
import { ReactComponent as WarningExclamation } from 'assets/big-warning-exclamation.svg';
import styled from '@emotion/styled';
import { colors } from 'Styles';
import { SoloButton } from './SoloButton';
import {
  SoloNegativeButton,
  SoloCancelButton,
  SoloButtonStyledComponent
} from 'Styles/CommonEmotions/button';

const maskStyle = {
  background: 'rgba(255, 255, 255, 0)'
};
const floatingStyle = {
  borderRadius: '10px'
};
const bodyStyle = {
  borderRadius: '10px'
};

const ContentContainer = styled.div`
  width: 250px;
  margin: 25px auto 50px;
  text-align: center;
`;

const WarningCircle = styled.div`
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 128px;
  height: 128px;
  border-radius: 100%;
  background: ${colors.flashlightGold};
  border: 2px solid ${colors.sunGold};
`;
const ContentText = styled.div`
  margin-top: 30px;
  font-size: 22px;
  color: ${colors.novemberGrey};
  width: 100%;
`;

const ButtonGroup = styled.div`
  display: flex;
  margin-top: 15px;
  justify-content: center;

  > button {
    min-width: 0;

    &:first-of-type {
      margin-right: 10px;
    }
  }
`;

interface Props extends ModalProps {
  confirmationTopic?: string;
  confirmText?: string;
  goForIt: () => any;
  cancel: () => any;
  visible?: boolean;
  isNegative?: boolean;
}

export const ConfirmationModal = (props: Props) => {
  const closeModal = (): void => {
    props.cancel();
  };

  const { confirmationTopic, confirmText, goForIt, isNegative } = props;

  return (
    <>
      <Modal
        visible={props.visible}
        footer={null}
        onCancel={closeModal}
        width={360}
        maskStyle={maskStyle}
        style={floatingStyle}
        bodyStyle={bodyStyle}>
        <ContentContainer>
          <WarningCircle>
            <WarningExclamation />
          </WarningCircle>
          <ContentText>
            Are you sure you want to{' '}
            {!!confirmationTopic ? confirmationTopic : 'remove this'}?
          </ContentText>

          <ButtonGroup>
            {isNegative ? (
              <SoloNegativeButton onClick={goForIt}>
                {!!confirmText ? confirmText : 'Confirm'}
              </SoloNegativeButton>
            ) : (
              <SoloButtonStyledComponent onClick={goForIt}>
                {!!confirmText ? confirmText : 'Confirm'}
              </SoloButtonStyledComponent>
            )}
            <SoloCancelButton onClick={closeModal}>Cancel</SoloCancelButton>
          </ButtonGroup>
        </ContentContainer>
      </Modal>
    </>
  );
};
