import React from 'react';
import ReactDom from 'react-dom';

import { App } from './App';

window.addEventListener('load', () => {
    ReactDom.render(<App />, document.getElementById('react-root'));
});
