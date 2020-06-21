import React, { forwardRef, MutableRefObject } from 'react';

import AddBox from '@material-ui/icons/AddBox';
import ArrowDownward from '@material-ui/icons/ArrowDownward';
import Check from '@material-ui/icons/Check';
import ChevronLeft from '@material-ui/icons/ChevronLeft';
import ChevronRight from '@material-ui/icons/ChevronRight';
import Clear from '@material-ui/icons/Clear';
import DeleteOutline from '@material-ui/icons/DeleteOutline';
import Edit from '@material-ui/icons/Edit';
import FilterList from '@material-ui/icons/FilterList';
import FirstPage from '@material-ui/icons/FirstPage';
import LastPage from '@material-ui/icons/LastPage';
import Remove from '@material-ui/icons/Remove';
import SaveAlt from '@material-ui/icons/SaveAlt';
import Search from '@material-ui/icons/Search';
import ViewColumn from '@material-ui/icons/ViewColumn';

export const tableIcons = {
  Add: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <AddBox {...props} ref={ref} />
  ),
  Check: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <Check {...props} ref={ref} />
  ),
  Clear: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <Clear {...props} ref={ref} />
  ),
  Delete: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <DeleteOutline {...props} ref={ref} />
  ),
  DetailPanel: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <ChevronRight {...props} ref={ref} />
  ),
  Edit: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <Edit {...props} ref={ref} />
  ),
  Export: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <SaveAlt {...props} ref={ref} />
  ),
  Filter: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <FilterList {...props} ref={ref} />
  ),
  FirstPage: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <FirstPage {...props} ref={ref} />
  ),
  LastPage: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <LastPage {...props} ref={ref} />
  ),
  NextPage: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <ChevronRight {...props} ref={ref} />
  ),
  PreviousPage: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <ChevronLeft {...props} ref={ref} />
  ),
  ResetSearch: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <Clear {...props} ref={ref} />
  ),
  Search: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <Search {...props} ref={ref} />
  ),
  SortArrow: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <ArrowDownward {...props} ref={ref} />
  ),
  ThirdStateCheck: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <Remove {...props} ref={ref} />
  ),
  ViewColumn: forwardRef(
    (
      props,
      ref:
        | ((instance: SVGSVGElement | null) => void)
        | MutableRefObject<SVGSVGElement | null>
        | null
    ) => <ViewColumn {...props} ref={ref} />
  )
};
