import React from 'react';
import styled from '@emotion/styled';
import { ReactComponent as SoloIcon } from 'assets/solo-logo-dark-text.svg';
import { colors } from 'Styles/colors';

const FooterContainer = styled.div`
  display: flex;
  min-height: 100%;
  flex-direction: row;
  justify-content: space-between;
  background: ${colors.februaryGrey};
`;

const Copyright = styled.div`
  color: ${colors.juneGrey};
  font-size: 16px;
  margin: 20px 90px;
`;

const Tagline = styled.div`
  align-content: flex-start;
  font-size: 10px;
  text-align: left;
`;

const IconContainer = styled.a`
  display: flex;
  flex-direction: row;
  margin: 10px 90px;
  filter: grayscale(100%) opacity(50%);
`;

const PrivacyPolicyContainer = styled.a`
  color: ${colors.juneGrey};
  font-weight: bold;
`;

export const Footer = () => {
  return (
    <FooterContainer>
      <Copyright>
        © {new Date().getFullYear()} solo.io, Inc. All Rights Reserved. |{' '}
        <PrivacyPolicyContainer
          href='https://www.solo.io/privacy-policy'
          target='_blank'>
          privacy policy
        </PrivacyPolicyContainer>
      </Copyright>
      <IconContainer href='https://www.solo.io' target='_blank'>
        <Tagline>Powered by</Tagline>
        <SoloIcon width='86' height='40' viewBox='0 0 300 100' />
      </IconContainer>
    </FooterContainer>
  );
};
