import * as fs from 'fs';
import * as path from 'path';
import { SundayApiSpec, SundayEndpoint } from '../universal-parser';

export function emitTS(spec: SundayApiSpec, outputDir: string) {
  const generatedDir = path.join(outputDir, 'generated');
  if (!fs.existsSync(generatedDir)) {
    fs.mkdirSync(generatedDir, { recursive: true });
  }

  for (const endpoint of spec.endpoints) {
    const slug = toKebabCase(endpoint.id);
    const className = toPascalCase(slug) + 'Client';
    const fileName = `${slug}.client.ts`;
    
    const content = `import { BaseApiClient } from '../../core/base-api.client';
import { z } from 'zod';

// TODO: Generate Zod schemas based on spec.schemas
export const ${toPascalCase(slug)}ResponseSchema = z.any();

export class ${className} extends BaseApiClient {
  /**
   * ${endpoint.summary || endpoint.id}
   * ${endpoint.description || ''}
   */
  async ${toCamelCase(slug)}(ctx: any, params?: any): Promise<any> {
    return this.${endpoint.method.toLowerCase()}(ctx, '${endpoint.path}', params);
  }
}
`;
    fs.writeFileSync(path.join(generatedDir, fileName), content);
  }
}

function toKebabCase(str: string): string {
  return str
    .replace(/([a-z])([A-Z])/g, '$1-$2')
    .replace(/[\s_]+/g, '-')
    .toLowerCase();
}

function toPascalCase(str: string): string {
  return str
    .split(/[-_ ]+/)
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join('');
}

function toCamelCase(str: string): string {
  const pascal = toPascalCase(str);
  return pascal.charAt(0).toLowerCase() + pascal.slice(1);
}
