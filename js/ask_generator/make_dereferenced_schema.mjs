//import { dereferenceSync } from 'dereference-json-schema';
import pkg from 'dereference-json-schema';
const { dereferenceSync } = pkg;
import { readFileSync, writeFileSync } from 'fs';

const schema = JSON.parse(readFileSync('../blueprint.schema.json', 'utf8'));
const dereferenced = dereferenceSync(schema);

// write to file
writeFileSync('../blueprint.schema.dereferenced.json', JSON.stringify(dereferenced, null, 2));

// to run this from bash, after 'npm i dereference-json-schema':
// node make_dereferenced_schema.mjs
