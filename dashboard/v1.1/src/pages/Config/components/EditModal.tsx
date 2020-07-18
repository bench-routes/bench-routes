import React, { Suspense, useState } from 'react';
import {
  AppBar,
  Button,
  Chip,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Tab,
  Tabs
} from '@material-ui/core';
import { TabPanel } from '../../Dashboard/SystemMetrics';
import Input from '../../Input/Input';
import { routeEntryType } from '../../../utils/types';

interface EditModalProps {
  isOpen: boolean;
  setOpen: (open: boolean) => void;
  selectedRoute: routeEntryType;
  updateConfigRoutes: (route: any) => void;
}

const EditModal = (props: EditModalProps) => {
  const { isOpen, setOpen, selectedRoute, updateConfigRoutes } = props;
  const [value, setValue] = useState(0);
  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const tabProps = (index: number) => {
    return {
      id: `simple-tab-${index}`,
      'aria-controls': `simple-tabpanel-${index}`
    };
  };

  const updateCurrModal = routes => {
    handleClose();
    updateConfigRoutes(routes);
  };

  return (
    <div>
      <Suspense fallback={<CircularProgress disableShrink />}>
        <Dialog
          fullWidth
          maxWidth="md"
          open={isOpen}
          onClose={handleClose}
          aria-labelledby="form-dialog-title"
        >
          <DialogTitle id="form-dialog-title">
            {selectedRoute ? (
              <div>
                {
                  <Chip
                    variant="outlined"
                    size="small"
                    label={selectedRoute.route}
                    clickable
                    color="primary"
                  />
                }
              </div>
            ) : (
              <div> </div>
            )}
          </DialogTitle>
          <DialogContent>
            <AppBar position="static">
              <Tabs
                value={value}
                onChange={handleChange}
                aria-label="Edit Route"
              >
                {Object.keys(selectedRoute).length !== 0 ? (
                  selectedRoute?.options?.map((options, index) => {
                    return (
                      <Tab
                        key={index}
                        label={options.method}
                        {...tabProps(index)}
                      />
                    );
                  })
                ) : (
                  <div> </div>
                )}
              </Tabs>
            </AppBar>
            {Object.keys(selectedRoute).length !== 0 ? (
              selectedRoute?.options?.map((options, index) => {
                return (
                  <TabPanel key={index} value={value} index={index}>
                    <Input
                      method={options.method}
                      headers={options.headers}
                      body={options.body}
                      params={options.params}
                      route={selectedRoute.route}
                      updateCurrModal={routes => updateCurrModal(routes)}
                      screenType="configScreen"
                    />
                  </TabPanel>
                );
              })
            ) : (
              <div> </div>
            )}
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">
              Close
            </Button>
          </DialogActions>
        </Dialog>
      </Suspense>
    </div>
  );
};

export default EditModal;
