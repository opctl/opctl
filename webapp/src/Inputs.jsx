import React from 'react';
import Form from "react-jsonschema-form";

export default function Inputs({
                                 value,
                               }) {

  const schema = {
    type: 'object',
    properties: {}
  };

  const uiSchema = {};
  if (value) {
    Object.entries(value).forEach(([name, param]) => {
      // type specific mapping
      if (param.array) {
        schema.properties[name] = Object.assign(
          {
            default: param.array.default,
            type: 'array',
          },
          param.array.constraints
        );

        uiSchema[name] = {
          'ui:description': param.array.description,
        };
      } else if (param.dir) {
        schema.properties[name] = {
          format: 'data-url',
          type: 'string',
        };

        uiSchema[name] = {'ui:description': param.dir.description};
      } else if (param.file) {
        schema.properties[name] = {
          format: 'data-url',
          type: 'string',
        };

        uiSchema[name] = {'ui:description': param.file.description};
      } else if (param.number) {
        schema.properties[name] = Object.assign(
          {
            default: param.number.default,
            type: 'number',
          },
          param.number.constraints,
        );

        uiSchema[name] = {
          'ui:description': param.number.description,
          "ui:widget": param.number.isSecret ? 'password' : 'text',
        };
      } else if (param.object) {
        schema.properties[name] = Object.assign(
          {
            default: param.object.default,
            type: 'object',
          },
          param.object.constraints,
        );

        uiSchema[name] = {
          'ui:description': param.object.description,
        };
      } else if (param.socket) {
        schema.properties[name] = {
          type: 'string',
          format: 'socket',
        };

        uiSchema[name] = {'ui:description': param.socket.description};
      } else if (param.string) {
        schema.properties[name] = Object.assign(
          {
            default: param.string.default,
            type: 'string',
          },
          param.string.constraints,
        );

        uiSchema[name] = {
          'ui:description': param.string.description,
          "ui:widget": param.string.isSecret ? 'password' : 'text',
        };
      }
    });
  }

  return (
    <div>
      <h2>Inputs:</h2>
      <Form schema={schema}
            uiSchema={uiSchema}
            liveValidate={true}
            onChange={() => console.log('changed')}
            onSubmit={() => console.log('submitted')}
            onError={() => console.log('erred')}/>
    </div>
  );
}
