// use dereferences schema to make a request to the passed host to do
// 'completion'.
//
// we will make a call to the local LM Studio instance (host passing required
// in case we need to go over the network).

import { readFileSync } from 'fs';
import pkg from 'dereference-json-schema';

// check number of args
if (process.argv.length < 3) {
    console.error('Usage: node build_furnace.mjs <host:port> [<model>]');
    console.error(' (useful host:port for LM Studio is localhost:5000)');
    process.exit(1);
}

const host = process.argv[2];
let model = false;
if (process.argv.length > 3) {
    model = process.argv[3];
}

// load schema
// it must not contain any $ref
let schema;
try {
    schema = JSON.parse(readFileSync('../../blueprint.schema.dereferenced.json', 'utf8'));
} catch (e) {
    // if file was not found, build a dereferenced schema from the original
    // schema
    const schemaOrig = JSON.parse(readFileSync('../../blueprint.schema.json', 'utf8'));
    // dereference it
    const { dereferenceSync } = pkg;
    schema = dereferenceSync(schemaOrig);
}
delete(schema['$schema']);
delete(schema['$id']);
schema['name'] = 'blueprint.schema.dereferenced.json';

// static prompt
var prompt = 'In Factorio, design a layout with stone furnaces filling up full belt of iron plates (that is, 15 items per second). Assume you have two full belts with one side of iron ore and other side of coal. You have yellow inserters available as well as small electric poles.'
prompt = 'Answer a prompt according to the blueprint schema. You MUST NOT reply with an empty response. Begin your response with `{"blueprint": {"description": "` and include the description of what you are building there. You will not reply with text and description, just plain JSON containing a root object named "blueprint", and then the response inside of it. Most of your response should be inside the "entitites" array, but you must populate the "item" with a "blueprint", and "icons" with a blueprint icon. Item prototype names you can use include: `stone-furnace`, `inserter`, `transport-belt`. x and y coordinates for inserters and belts are centered at integer + 0.5 because they are 1x1, while furnaces are 2x2 so they are centered at an integer position. Prototype name for small electric pole is `electric-pole`. Orientation is an integer 0-8. Root item of the response inside the object keyed as `blueprint` must have the prototype (i.e. property `item`) set to `blueprint`.\n\n' + prompt;

// build chat completion request.
// example:
// '{
//    "model": "your-model-name-here",
//    "messages": [
//      { "role": "system", "content": "" },
//      { "role": "assistant", "content": "prompt comes here" }
//    ],
//    "temperature": 0.7,
//    "max_tokens": -1,
//    "stream": false,
//    "response_format": {
//       "type": "json_schema",
//       "strict": true,
//       "json_schema": <dereferenced schema>
//    }
//  }'

let request = {
    messages: [
        { role: 'system', content: 'You are a helpful assistant and will return a response in JSON.' },
        { role: 'assistant', content: prompt }
    ],
    temperature: 0.7,
    max_tokens: -1,
    stream: false,
};
if (model) {
    request.model = model
};
if (true) {
    request.response_format = {
        type: 'json_schema',
        json_schema: {
            name: 'blueprint.dereferenced.schema.json',
            strict: true,
            schema: schema
        },
    }
}

// make the request
fetch(`http://${host}/v1/chat/completions`, {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify(request)
}).then(response => {
    if (response.ok) {
        return response.json();
    }
    throw new Error('Failed to get response: ' + response.statusText);
}).then(data => {
    console.log('Response:', data);
    let raw_json = data.choices[0].message.content;
    // we could decode the response, but we could also just write it out:
    console.log(raw_json);
}).catch(error => {
    console.error('Error:', error);
});

