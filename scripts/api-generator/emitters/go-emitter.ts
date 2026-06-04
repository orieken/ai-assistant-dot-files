import * as fs from 'fs';
import * as path from 'path';
import { SundayApiSpec, SundayEndpoint } from '../universal-parser';

export function emitGo(spec: SundayApiSpec, outputDir: string) {
  const generatedDir = path.join(outputDir, 'client', 'generated');
  if (!fs.existsSync(generatedDir)) {
    fs.mkdirSync(generatedDir, { recursive: true });
  }

  for (const endpoint of spec.endpoints) {
    const slug = toSnakeCase(endpoint.id);
    const structName = toPascalCase(slug) + 'ApiClient';
    const fileName = `${slug}_client.go`;
    
    const content = `package generated

import (
	"context"
	"github.com/orieken/go-sunday/core"
	"github.com/orieken/go-sunday/client"
)

// ${structName} handles ${endpoint.method} ${endpoint.path}
type ${structName} struct {
	*client.BaseApiClient
}

func New${structName}(adapter core.IHttpAdapter, opts ...client.ClientOption) *${structName} {
	return &${structName}{client.NewBaseApiClient(adapter, opts...)}
}

func (p *${structName}) ${toPascalCase(slug)}(ctx context.Context) (*core.HttpResponse[any], error) {
	return client.${toPascalCase(endpoint.method.toLowerCase())}[any](ctx, p.BaseApiClient, "${endpoint.path}")
}
`;
    fs.writeFileSync(path.join(generatedDir, fileName), content);
  }
}

function toSnakeCase(str: string): string {
  return str
    .replace(/([a-z])([A-Z])/g, '$1_$2')
    .replace(/[\s-]+/g, '_')
    .toLowerCase();
}

function toPascalCase(str: string): string {
  return str
    .split(/[-_ ]+/)
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join('');
}
