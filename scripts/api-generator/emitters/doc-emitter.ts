import * as fs from 'fs';
import * as path from 'path';
import { SundayApiSpec, SundayEndpoint } from '../universal-parser';

export function emitDocs(spec: SundayApiSpec, outputDir: string) {
  for (const endpoint of spec.endpoints) {
    const slug = toKebabCase(endpoint.id);
    const featureDir = path.join(outputDir, 'features', slug);
    if (!fs.existsSync(featureDir)) {
      fs.mkdirSync(featureDir, { recursive: true });
    }

    const analysisContent = generateAnalysis(endpoint);
    const archNotesContent = generateArchNotes(endpoint);

    fs.writeFileSync(path.join(featureDir, 'analysis.md'), analysisContent);
    fs.writeFileSync(path.join(featureDir, 'architecture-notes.md'), archNotesContent);
  }
}

function toKebabCase(str: string): string {
  return str
    .replace(/([a-z])([A-Z])/g, '$1-$2')
    .replace(/[\s_]+/g, '-')
    .toLowerCase();
}

function generateAnalysis(endpoint: SundayEndpoint): string {
  return `# Analysis: ${endpoint.summary || endpoint.id}

## Requirements
- **Method**: ${endpoint.method}
- **Path**: \`${endpoint.path}\`
- **Description**: ${endpoint.description || 'No description provided.'}

## Parameters
${endpoint.parameters.map(p => `- **${p.name}** (${p.in}): ${p.description || ''}`).join('\n') || 'None'}

## Success Criteria
- [ ] Endpoint returns successful response for valid inputs.
- [ ] Appropriate error codes are returned for invalid inputs.
`;
}

function generateArchNotes(endpoint: SundayEndpoint): string {
  return `# Architecture Notes: ${endpoint.id}

## Interface Design
- **Service**: Generated API Client
- **Endpoint**: \`${endpoint.method} ${endpoint.path}\`

## Request Shape
${endpoint.requestBody ? 'Requires request body (see OpenAPI spec)' : 'No request body'}

## Response Shape
${Object.entries(endpoint.responses).map(([code, resp]) => `- **${code}**: ${resp.description}`).join('\n')}

## Security
- [ ] Authentication required
- [ ] Authorization roles verified
`;
}
