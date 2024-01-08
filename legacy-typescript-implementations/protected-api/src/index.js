const express = require('express');
const logger = require('morgan');
const jwksRsa = require('jwks-rsa');
var jwt = require('express-jwt');
const uuid = require('uuid');
const { Issuer } = require('openid-client');

const port = process.env.CLIENT_PORT || 5010;
const oidc_issuer_url = process.env.OIDC_ISSUER_URL;

const app = express();

const objects = {};

console.log('OIDC_ISSUER_URL', oidc_issuer_url);

objects[uuid.v4()] = {title: 'Test object'}

app.use(express.json());
app.use(logger('combined'));


// OIDC discovery to locate the JWKS URI used to validate JWTs
Issuer.discover(oidc_issuer_url)
    .then(function (issuer) {
        console.log('Discovered issuer %s %O', issuer.issuer, issuer.metadata);

        // Install JWT middleware using 'issuer.jwks_uri' and caching of keys
        app.use(jwt({
            secret: jwksRsa.expressJwtSecret({
                jwksUri: issuer.jwks_uri,
                cache: true,  // Enable JWT key cache
                timeout: 3600 // Key cache timeout, seconds
            }),
            algorithms: [ 'RS256' ],
            requestProperty: 'auth'
            // Here we could check more 'static' properties
            // audience: ...
            // issuer: ...
        }));

        app.get('/objects', (req, res) => {
            console.log('Read list of object IDs');
            res.send(Object.keys(objects));
        });

        app.post('/object',
                 (req, res) => {
                     const id = uuid.v4();
                     objects[id] = req.body.data;
                     console.log('Created new object with ID', id, 'data', objects[id]);
                     res.send({id});
        });

        app.use(function (err, req, res, next) {
            console.log('Error handler', err);
            if (err.name === 'UnauthorizedError') {
                res.status(401).send('invalid token');
            } else {
                res.status(err.status).send(err);
            }
        });

        app.listen(port, () => {
            console.log(`Object store listening on port ${port}!`);
        });
    });
