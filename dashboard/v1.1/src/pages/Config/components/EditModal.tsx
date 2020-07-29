import React, { Suspense, useState, useEffect } from 'react';
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
import { truncate } from '../../../utils/stringManipulations';

interface EditModalProps {
  isOpen: boolean;
  setOpen: (open: boolean) => void;
  selectedRoute: routeEntryType;
  updateConfigRoutes: (route: any) => void;
}

const EditModal = (props: EditModalProps) => {
  const { isOpen, setOpen, selectedRoute, updateConfigRoutes } = props;
  const [modalUrl, setModalUrl] = useState<string>('');
  const [value, setValue] = useState<number>(0);
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

  useEffect(() => {
    setValue(0);
    setModalUrl(selectedRoute.route);
  }, [selectedRoute]);

  const updateCurrentModal = (routes, URL: string) => {
    setModalUrl(URL);
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
                    size="medium"
                    style={{ fontSize: '14px' }}
                    label={truncate(modalUrl, 100)}
                    clickable
                    color="primary"
                  />
                }
              </div>
            ) : (
              <> </>
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
                  <> </>
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
                      labels={options.labels}
                      route={selectedRoute.route}
                      updateCurrentModal={(routes, URL) =>
                        updateCurrentModal(routes, URL)
                      }
                      screenType="config-screen"
                    />
                  </TabPanel>
                );
              })
            ) : (
              <> </>
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
