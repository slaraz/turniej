//w katalogu /proto
//grpc_tools_node_protoc --js_out=import_style=commonjs,binary:../gra_nodejs/src --grpc_out=generate_package_definition:../gra_nodejs/src gra.proto

const PROTO_PATH = `../proto/gra.proto`;
const SERWER_URL = `localhost:50051`;

const grpc = require(`@grpc/grpc-js`);
const graMessages = require(`./gra_pb.js`);
const graServices = require(`./gra_grpc_pb.js`);

const packageDefinition = grpc.loadPackageDefinition(graServices);
const meta = new grpc.Metadata();

konfiguracjaGry = new graMessages.KonfiguracjaGry();
konfiguracjaGry.setLiczbagraczy(2);
  
const client = new packageDefinition.Gra(SERWER_URL, grpc.credentials.createInsecure())

client.nowyMecz(konfiguracjaGry, (err, res) => {
  console.log(err)
  console.log(res.getGraid())
});

