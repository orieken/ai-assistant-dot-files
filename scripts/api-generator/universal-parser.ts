import axios from 'axios';
import jsYaml from 'js-yaml';
import { OpenAPIV3 } from 'openapi-types';

export interface SundayEndpoint {
  id: string; // operationId or slug
  path: string;
  method: string;
  summary?: string;
  description?: string;
  parameters: OpenAPIV3.ParameterObject[];
  requestBody?: OpenAPIV3.RequestBodyObject;
  responses: Record<string, OpenAPIV3.ResponseObject>;
}

export interface SundayApiSpec {
  info: OpenAPIV3.InfoObject;
  endpoints: SundayEndpoint[];
  schemas: Record<string, OpenAPIV3.SchemaObject>;
}

export async function parseOpenApi(source: string): Promise<SundayApiSpec> {
  let data: any;
  if (source.startsWith('http')) {
    const response = await axios.get(source);
    data = response.data;
  } else {
    const fs = await import('fs');
    const content = fs.readFileSync(source, 'utf8');
    data = content;
  }
  
  const rawSpec = typeof data === 'string' ? jsYaml.load(data) as any : data;
  const spec = rawSpec as OpenAPIV3.Document;

  if (!spec.openapi || !spec.openapi.startsWith('3.')) {
    throw new Error('Only OpenAPI 3.x is supported currently.');
  }

  const endpoints: SundayEndpoint[] = [];
  
  for (const [path, pathItem] of Object.entries(spec.paths)) {
    if (!pathItem) continue;

    for (const method of ['get', 'post', 'put', 'delete', 'patch']) {
      const operation = (pathItem as any)[method] as OpenAPIV3.OperationObject;
      if (!operation) continue;

      const id = operation.operationId || `${method}-${path.replace(/\//g, '-').replace(/[{}]/g, '')}`.replace(/-+/g, '-').replace(/^-|-$/g, '');
      
      endpoints.push({
        id,
        path,
        method: method.toUpperCase(),
        summary: operation.summary,
        description: operation.description,
        parameters: (operation.parameters || []) as OpenAPIV3.ParameterObject[],
        requestBody: operation.requestBody as OpenAPIV3.RequestBodyObject,
        responses: operation.responses as Record<string, OpenAPIV3.ResponseObject>,
      });
    }
  }

  return {
    info: spec.info,
    endpoints,
    schemas: (spec.components?.schemas || {}) as Record<string, OpenAPIV3.SchemaObject>,
  };
}
