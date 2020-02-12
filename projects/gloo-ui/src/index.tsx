import * as React from 'react';
import * as ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import './fontFace.css';
import { GlooIApp } from './GlooIApp';
import * as serviceWorker from './serviceWorker';
import { globalStore } from './store';
import { ErrorBoundary } from 'Components/Features/Errors/ErrorBoundary';
import { SWRConfig } from 'swr';
import { SoloWarning } from 'Components/Common/SoloWarningContent';
import { BrowserRouter } from 'react-router-dom';

ReactDOM.render(
  <ErrorBoundary fallback={<div> there was an error</div>}>
    <Provider store={globalStore}>
      <BrowserRouter>
        <GlooIApp />
      </BrowserRouter>
    </Provider>
  </ErrorBoundary>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
