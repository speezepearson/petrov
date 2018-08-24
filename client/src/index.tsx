import React from 'react';
import ReactDom from 'react-dom';

window.addEventListener('load', () => {
    ReactDom.render(<div>Hello world!</div>, document.getElementById('react-root'));
})
