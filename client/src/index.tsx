import React from 'react';
import ReactDom from 'react-dom';
import jquery from 'jquery';
import { App } from './App';

alert("Some browsers block sounds unless you've interacted with the page. Click Continue to interact with the page.");

window.addEventListener('load', () => {
    const hrefComponents = window.location.href.replace(/\/+$/, '').split('/');
    let password = hrefComponents[hrefComponents.length - 1];
    ReactDom.render(
        <App password={password} />,
        document.getElementById('react-root'),
    );
});
