import { combineEpics } from 'redux-observable';
import { browserHistory } from 'react-router';
import { successOf } from '@/services/api/ajaxEpic';
import { getBoardSelector } from '@/Home/Board/duck';
import { DELETE } from '@/Home/Board/List/Card/duck';

const closeCardModalOnDeleteSuccessEpic = (action$, store) => action$.ofType(successOf(DELETE))
  .do(() => {
    const board = getBoardSelector(store.getState());
    browserHistory.push(`/${board.owner.login}/${board.slug}`);
  })
  .ignoreElements();

export const epics = combineEpics(/* eslint import/prefer-default-export: 0 */
  closeCardModalOnDeleteSuccessEpic,
);
