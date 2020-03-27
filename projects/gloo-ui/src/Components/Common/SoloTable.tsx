import styled from '@emotion/styled';
import Table from 'antd/lib/table';
import * as React from 'react';
import { soloConstants } from 'Styles';
import { colors, hslToHSLA } from 'Styles/colors';

// To restyle table to match spec later
export const TableContainer = styled.div`
  box-shadow: 0px 4px 9px ${colors.boxShadow};

  .ant-table-wrapper {
    background: ${hslToHSLA(colors.marchGrey, 0.15)};

    .ant-table-thead {
      border-radius: ${soloConstants.radius}px ${soloConstants.radius}px 0 0;

      tr {
        background: ${colors.marchGrey};

        .ant-table-column-title {
          cursor: default;
          font-weight: 600;
          color: ${colors.novemberGrey};
        }
      }
    }

    .ant-table-tbody {
      background: white;

      .ant-table-row {
        position: relative;

        > td {
          border-color: ${colors.februaryGrey};
        }
        &:hover {
          > td {
            background: ${hslToHSLA(colors.marchGrey, 0.5)};
          }
        }
      }
    }

    .ant-table-pagination {
      &.ant-pagination {
        margin: 0;
      }

      &[unselectable='unselectable'] {
        opacity: 0.25;
      }

      .ant-pagination-prev,
      .ant-pagination-next {
      }

      a,
      a.ant-pagination-item-link,
      .ant-pagination-item-active {
        background: none;
        border: none;
        color: ${colors.septemberGrey};
        line-height: 30px;
        height: 30px;
      }
    }

    .ant-empty-description {
      color: white;

      &::after {
        content: 'No Matches';
        position: absolute;
        left: 0;
        right: 0;
        text-align: center;
        color: ${colors.juneGrey};
      }
    }
  }
`;

export interface TableProps {
  columns: any[];
  dataSource: any[];
  formComponent?: React.FC;
  title?: string;
}

// TODO: figure out if edit row should always be shown or always be last row
const EditableRow = ({ lastRowID, formComponent, isEmpty, ...props }: any) => {
  const isLastRow = lastRowID === props['data-row-key'];
  const FormComponent = formComponent;

  return (
    <>
      {isLastRow && !!formComponent ? (
        <tr>
          <FormComponent />
        </tr>
      ) : (
        <tr {...props} />
      )}
    </>
  );
};

export const SoloTable = (props: TableProps) => {
  const components = {
    body: {
      row: EditableRow
    }
  };

  const lastRowID =
    props.dataSource.length > 0
      ? props.dataSource[props.dataSource.length - 1].key
      : true;

  return (
    <TableContainer>
      <Table
        title={props.title ? () => <b>{props.title}</b> : undefined}
        dataSource={props.dataSource}
        columns={props.columns}
        components={components}
        pagination={{ defaultPageSize: 20 }}
        onRow={record => {
          return {
            ...record,
            lastRowID,
            isEmpty: props.dataSource.length === 1,
            cols: props.columns,
            formComponent: props.formComponent
          };
        }}
      />
    </TableContainer>
  );
};
