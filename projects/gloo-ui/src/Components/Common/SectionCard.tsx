import styled from '@emotion/styled';
import * as React from 'react';
import { colors, soloConstants } from 'Styles';
import { CardCSS } from 'Styles/CommonEmotions/card';
import { HealthIndicator } from './HealthIndicator';
import { Status } from 'proto/solo-kit/api/v1/status_pb';

const CardBlock = styled.div`
  ${CardCSS};
  margin-bottom: 30px;
  padding: 0;
  @media (max-width: 1380px) {
    margin-bottom: 45px;
  }
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  width: 100%;
  background: ${colors.marchGrey};
  padding: 13px;
  border-radius: ${soloConstants.radius}px ${soloConstants.radius}px 0 0;
`;

const HeaderImageHolder = styled.div`
  margin-right: 15px;
  height: 33px;
  width: 33px;
  border-radius: 100%;
  background: white;
  display: flex;
  justify-content: center;
  align-items: center;

  img,
  svg {
    width: 30px;
    max-height: 30px;
  }
`;

const HeaderTitleSection = styled.div`
  max-width: calc(100% - 300px);
`;
const HeaderTitleName = styled.div`
  width: 100%;
  font-size: 22px;
  color: ${colors.novemberGrey};
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
`;

const SecondaryInformation = styled.div`
  display: flex;
  align-items: center;
`;
const SecondaryInformationSection = styled.div`
  font-size: 14px;
  line-height: 22px;
  height: 22px;
  padding: 0 12px;
  color: ${colors.novemberGrey};
  background: white;
  margin-left: 13px;
  border-radius: ${soloConstants.largeRadius}px;
`;
const SecondaryInformationTitle = styled.span`
  font-weight: bold;
`;

const HealthContainer = styled.div`
  display: flex;
  justify-content: flex-end;
  align-items: center;
  flex: 1;
  text-align: right;
  font-size: 16px;
  font-weight: 600;
  color: ${colors.novemberGrey};
`;

const CloseIcon = styled.div`
  font-size: 21px;
  line-height: 17px;
  margin-left: ${soloConstants.largeBuffer}px;
  margin-top: 2px;
  font-weight: 100;
  color: ${colors.juneGrey};
  cursor: pointer;
`;

type BodyContainerProps = { noPadding: boolean };
const BodyContainer = styled.div`
  padding: ${(props: BodyContainerProps) => (props.noPadding ? '' : '20px;')};
`;

interface Props {
  cardName: string;
  logoIcon?: React.ReactNode;
  headerSecondaryInformation?: {
    title?: string;
    value: string;
  }[];
  health?: Status.StateMap[keyof Status.StateMap];
  healthMessage?: string;
  onClose?: () => any;
  secondaryComponent?: React.ReactNode;
  noPadding?: boolean;
}

export const SectionCard: React.FunctionComponent<Props> = props => {
  const {
    logoIcon,
    cardName,
    children,
    headerSecondaryInformation,
    health,
    healthMessage,
    onClose,
    secondaryComponent,
    noPadding
  } = props;

  return (
    <CardBlock>
      <Header>
        {logoIcon && <HeaderImageHolder>{logoIcon}</HeaderImageHolder>}
        <HeaderTitleSection>
          <HeaderTitleName>{cardName}</HeaderTitleName>
        </HeaderTitleSection>
        {!!secondaryComponent && (
          <SecondaryInformation>
            <SecondaryInformationSection>
              {secondaryComponent}
            </SecondaryInformationSection>
          </SecondaryInformation>
        )}
        {!!headerSecondaryInformation && (
          <SecondaryInformation>
            {headerSecondaryInformation.map(info => {
              return (
                <SecondaryInformationSection key={info.value}>
                  {!!info.title && (
                    <SecondaryInformationTitle>
                      {info.title}:{' '}
                    </SecondaryInformationTitle>
                  )}
                  {info.value}
                </SecondaryInformationSection>
              );
            })}
          </SecondaryInformation>
        )}
        {(!!health || health === 0) && (
          <HealthContainer>
            {healthMessage || ''}
            <HealthIndicator healthStatus={health} />
          </HealthContainer>
        )}
        {onClose && <CloseIcon onClick={onClose}>X</CloseIcon>}
      </Header>
      <BodyContainer noPadding={noPadding}>{children}</BodyContainer>
    </CardBlock>
  );
};
