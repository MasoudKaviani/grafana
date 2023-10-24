import { css } from '@emotion/css';
import React, { useState } from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { Dropdown, useStyles2 } from '@grafana/ui';
import { Box } from '@grafana/ui/src/unstable';
import { DashNavButton } from 'app/features/dashboard/components/DashNav/DashNavButton';

export interface Props {}

export function PageHistoryButton(props: Props) {
  const styles = useStyles2(getStyles);
  const [isOpen, setIsOpen] = useState(false);

  const renderPopover = () => {
    return (
      /* eslint-disable-next-line jsx-a11y/no-static-element-interactions, jsx-a11y/click-events-have-key-events */
      <div className={styles.popover} onClick={(evt) => evt.stopPropagation()}>
        <div className={styles.heading}>Page history</div>
      </div>
    );
  };

  return (
    <Dropdown overlay={renderPopover} placement="bottom" onVisibleChange={setIsOpen}>
      <Box marginLeft={1}>
        <DashNavButton icon="clock-nine" tooltip="Settings" />
      </Box>
    </Dropdown>
  );
}

const getStyles = (theme: GrafanaTheme2) => {
  return {
    popover: css({
      display: 'flex',
      padding: theme.spacing(2),
      flexDirection: 'column',
      background: theme.colors.background.primary,
      boxShadow: theme.shadows.z3,
      borderRadius: theme.shape.borderRadius(),
      border: `1px solid ${theme.colors.border.weak}`,
      zIndex: 1,
      marginRight: theme.spacing(2),
    }),
    heading: css({
      fontWeight: theme.typography.fontWeightMedium,
      paddingBottom: theme.spacing(2),
    }),
    options: css({
      display: 'grid',
      gridTemplateColumns: '1fr 50px',
      rowGap: theme.spacing(1),
      columnGap: theme.spacing(2),
    }),
  };
};
