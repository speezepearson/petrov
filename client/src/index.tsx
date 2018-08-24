import React from 'react';
import ReactDom from 'react-dom';

import { App } from './App';

window.addEventListener('load', () => {
    ReactDom.render(
        <App playerName="Alice" />,
        document.getElementById('react-root'),
    );
});
