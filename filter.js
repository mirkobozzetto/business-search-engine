// filter.js  ←  node filter.js 62020
const fs  = require("fs");
const csv = require("csv-parser");

const DIR  = "bce_mai_2025";
const CODE = process.argv[2];

if (!CODE || !/^\d{5}$/.test(CODE)) {
  console.error("usage : node filter.js <NACE 5 chiffres>");
  process.exit(1);
}

const OUT    = `filtre_${CODE}.csv`;
const header = `"EntityNumber","ActivityGroup","NaceVersion","NaceCode","Classification"\n`;
const out    = fs.createWriteStream(OUT, "utf8");
out.write(header);

let n = 0;

fs.createReadStream(`${DIR}/activity.csv`)
  .pipe(csv({ separator: ",", quote: '"' }))
  .on("data", (r) => {
    if (r.NaceCode === CODE) {
      n++;
      out.write(
        `"${r.EntityNumber}","${r.ActivityGroup || ""}","${r.NaceVersion || ""}","${r.NaceCode}","${r.Classification || ""}"\n`
      );
    }
  })
  .on("end", () => {
    out.end();
    console.log(`✅ ${OUT} créé (${n} lignes)`);
  });
