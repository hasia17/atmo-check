const fs = require('fs');
const svg = fs.readFileSync('./public/poland.svg', 'utf8');

const voivodeships = [
  'dolnoslaskie','kujawsko-pomorskie','lubelskie','lubuskie','lodzkie',
  'malopolskie','mazowieckie','opolskie','podkarpackie','podlaskie',
  'pomorskie','slaskie','swietokrzyskie','warminsko-mazurskie',
  'wielkopolskie','zachodniopomorskie'
];

// Tokenize path data into commands and numbers
function tokenize(d) {
  return d.match(/[MmLlHhVvCcSsQqTtAaZz]|[-+]?(?:\d+\.?\d*|\.\d+)(?:[eE][-+]?\d+)?/g) || [];
}

// Parse SVG path into absolute coordinate points, split by subpath
function parseSubpaths(d) {
  const tokens = tokenize(d);
  let x = 0, y = 0, startX = 0, startY = 0;
  const subpaths = [];
  let pts = [];
  let cmd = '';
  let i = 0;

  function num() { return parseFloat(tokens[i++]); }

  while (i < tokens.length) {
    const t = tokens[i];
    if (/^[MmLlHhVvCcSsQqTtAaZz]$/.test(t)) {
      cmd = t;
      i++;
      if (cmd === 'Z' || cmd === 'z') {
        if (pts.length) subpaths.push(pts);
        pts = [];
        x = startX; y = startY;
      }
      continue;
    }

    switch (cmd) {
      case 'M': x=num(); y=num(); pts=[{x,y}]; startX=x; startY=y; cmd='L'; break;
      case 'm': {
        if (pts.length) subpaths.push(pts);
        pts = [];
        x+=num(); y+=num(); pts.push({x,y}); startX=x; startY=y; cmd='l'; break;
      }
      case 'L': x=num(); y=num(); pts.push({x,y}); break;
      case 'l': x+=num(); y+=num(); pts.push({x,y}); break;
      case 'H': x=num(); pts.push({x,y}); break;
      case 'h': x+=num(); pts.push({x,y}); break;
      case 'V': y=num(); pts.push({x,y}); break;
      case 'v': y+=num(); pts.push({x,y}); break;
      case 'C': { num();num();num();num(); x=num(); y=num(); pts.push({x,y}); break; }
      case 'c': { num();num();num();num(); x+=num(); y+=num(); pts.push({x,y}); break; }
      case 'S': { num();num(); x=num(); y=num(); pts.push({x,y}); break; }
      case 's': { num();num(); x+=num(); y+=num(); pts.push({x,y}); break; }
      case 'Q': { num();num(); x=num(); y=num(); pts.push({x,y}); break; }
      case 'q': { num();num(); x+=num(); y+=num(); pts.push({x,y}); break; }
      case 'T': { x=num(); y=num(); pts.push({x,y}); break; }
      case 't': { x+=num(); y+=num(); pts.push({x,y}); break; }
      case 'A': { num();num();num();num();num(); x=num(); y=num(); pts.push({x,y}); break; }
      case 'a': { num();num();num();num();num(); x+=num(); y+=num(); pts.push({x,y}); break; }
      default: i++; break; // skip unknown
    }
  }
  if (pts.length) subpaths.push(pts);
  return subpaths;
}

function bboxCenter(points) {
  const xs = points.map(p=>p.x), ys = points.map(p=>p.y);
  const cx = Math.round((Math.min(...xs)+Math.max(...xs))/2);
  const cy = Math.round((Math.min(...ys)+Math.max(...ys))/2);
  const area = (Math.max(...xs)-Math.min(...xs))*(Math.max(...ys)-Math.min(...ys));
  return {cx, cy, area};
}

for (const id of voivodeships) {
  let match = svg.match(new RegExp(`<path[^>]*id="${id}"[^>]*d="([^"]+)"`, 's'));
  if (!match) match = svg.match(new RegExp(`<path[^>]*d="([^"]+)"[^>]*id="${id}"`, 's'));
  if (!match) { console.log(`  /* ${id}: NOT FOUND */`); continue; }

  const subpaths = parseSubpaths(match[1]);
  // Use largest subpath (by bounding box area) — ignores tiny island dots
  let best = null;
  for (const sp of subpaths) {
    if (sp.length < 4) continue;
    const b = bboxCenter(sp);
    if (!best || b.area > best.area) best = {...b, n: sp.length};
  }

  if (!best) {
    const all = subpaths.flat();
    best = all.length ? {...bboxCenter(all), n: all.length} : {cx:0,cy:0,area:0,n:0};
  }

  console.log(`  { id: '${id}', cx: ${best.cx}, cy: ${best.cy} },`);
}
