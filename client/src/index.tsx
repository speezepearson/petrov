import React from 'react';
import ReactDom from 'react-dom';
import jquery from 'jquery';
import { App } from './App';

window.addEventListener('load', () => {
    const hrefComponents = window.location.href.replace(/\/+$/, '').split('/');
    let password = hrefComponents[hrefComponents.length - 1];
    ReactDom.render(
        <App password={password} />,
        document.getElementById('react-root'),
    );
});
