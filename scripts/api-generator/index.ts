import * as path from 'path';
import { parseOpenApi } from './universal-parser';
import { emitDocs } from './emitters/doc-emitter';
import { emitTS } from './emitters/ts-emitter';
import { emitGo } from './emitters/go-emitter';

async function main() {
  const args = process.argv.slice(2);
  const url = args[0];
  
  if (!url) {
    console.error('Usage: ts-node index.ts <swagger-url>');
    process.exit(1);
  }

  try {
    console.log(`🚀 Ingesting API from: ${url}`);
    const spec = await parseOpenApi(url);
    
    const rootDir = path.resolve(path.dirname(new URL(import.meta.url).pathname), '../../');
    const goDir = '/Users/oscarrieken/Projects/Rieken/go-sunday';
    const tsDir = '/Users/oscarrieken/Projects/Rieken/sunday-monorepo/packages/clients';

    console.log(`rootDir: ${rootDir}`);
    console.log(`Docs output: ${path.join(rootDir, 'docs')}`);

    console.log('📝 Generating Docs...');
    emitDocs(spec, path.join(rootDir, 'docs'));

    console.log('TypeScript: Generating Clients...');
    emitTS(spec, tsDir);

    console.log('Go: Generating Clients...');
    emitGo(spec, goDir);

    console.log('✅ Ingestion complete!');
  } catch (error) {
    console.error('❌ Error during ingestion:', error);
    process.exit(1);
  }
}

main();
