import { browserHistory } from 'react-router';
import { combineEpics } from 'redux-observable';
import { Observable } from 'rxjs/Observable';
import { successOf } from '../../../services/api/ajaxEpic';
import { UPDATE, DELETE } from '../duck';
import { getBoards } from '../../Menu/BoardList/duck';
import { hideModal } from '../../../components/Modal/duck';

// types
export const MODAL_DELETE_BOARD = 'pulpe/home/board/header/modalDeleteBoard';

// epics
const redirectOnBoardDeletionEpic = (action$, store) => action$.ofType(successOf(DELETE))
  .do(() =>
    setTimeout(() => {
      const boards = getBoards(store.getState());
      if (boards && boards.length > 0) {
        browserHistory.push(`/${boards[0].slug}`);
      } else {
        browserHistory.push('/');
      }
    }, 450)
  )
  .ignoreElements();

const closeModalOnDeletionEpic = action$ => action$.ofType(successOf(DELETE))
  .map(hideModal);

const redirectOnBoardUpdateEpic = action$ => action$.ofType(successOf(UPDATE))
  .do((action) => browserHistory.push(`/${action.response.entities.boards[action.response.result].slug}`))
  .mergeMap(() => Observable.empty());

export const epics = combineEpics(
  redirectOnBoardDeletionEpic,
  closeModalOnDeletionEpic,
  redirectOnBoardUpdateEpic
);
