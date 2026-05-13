import fs from 'fs';
import path from 'path';
import express from 'express';
import spdy from 'spdy';

const projectRoot = path.resolve(__dirname, '..');

let sslKeyPath = process.env.SSL_KEY_PATH || '.keys/server-key.pem';
let sslCertPath = process.env.SSL_CERT_PATH || '.keys/server-cert.pem';
let clientPath = process.env.CLIENT_PATH || '../client/build';
const port = process.env.PORT || 8443;

sslKeyPath = path.resolve(projectRoot, sslKeyPath);
sslCertPath = path.resolve(projectRoot, sslCertPath);
clientPath = path.resolve(projectRoot, clientPath);

const app = express();

app.post('/api/login', (req, res) => {
  res.status(200).send('stub');
});

app.post('/api/logout', (req, res) => {
  res.status(200).send('stub');
});

app.get('/api/filebrowser', (req, res) => {
  res.status(200).send('stub');
});

// fallback route for paths starting in '/api'
app.use('/api', (req, res) => {
  res.status(404).send('Not found');
});

app.use(express.static(clientPath));

// fallback route for static content: return the contents of index.html rather
// than a 404. This is to support pushState-based routing on the client.
app.use((req, res, next) => {
  req.url = '/index.html';
  express.static(clientPath)(req, res, next);
});

spdy
  .createServer(
    {
      key: fs.readFileSync(sslKeyPath),
      cert: fs.readFileSync(sslCertPath),
    },
    app,
  )
  .listen(port, () => {
    console.log(`Listening on port ${port}`);
  });
