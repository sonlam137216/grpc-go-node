const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

const express = require('express');
const app = express();
const port = 5000;

const packageDefinition = protoLoader.loadSync(
  '../reader/protobuf/read-excel.proto',
  {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
  }
);

const exampleProto = grpc.loadPackageDefinition(packageDefinition);

const sayHelloGRPC = () => {
  return new Promise((resolve) => {
    const client = new exampleProto.ExampleService(
      'localhost:50051',
      grpc.credentials.createInsecure()
    );

    client.SayHello({ name: 'World' }, (err, response) => {
      if (err) {
        console.error(err);
      }
      resolve(response.rows);
    });
  });
};

app.get('/', async (req, res) => {
  const data = await sayHelloGRPC();
  data.forEach((item) => {
    console.log(item);
  });
  res.send({ status: 'done' });
});

app.listen(port, () => {
  console.log(`Server listening at ${port}`);
});
