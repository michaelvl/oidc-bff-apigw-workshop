const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const logger = require('morgan');

const port = process.env.CLIENT_PORT || 5000;

const app = express();
app.use(logger('combined'));

app.use('/login', createProxyMiddleware({
    target: "http://localhost:5010",
    changeOrigin: true  //,
    //pathRewrite: {'^/login' : ''},
}));

app.use('/api', createProxyMiddleware({
    target: "http://localhost:5020",
    changeOrigin: true,
    pathRewrite: {'^/api' : ''},
}));

app.use('/', createProxyMiddleware({
    target: "http://localhost:5030",
    changeOrigin: true,
}));

app.listen(port, () => {
    console.log(`Router listening on port ${port}!`);
});
