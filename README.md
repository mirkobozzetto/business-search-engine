## Overview

This project is a **minimal Node.js CLI** that extracts any **NACE-BEL 5-digit sector** from the Belgian Crossroads Bank for Enterprises (CBE) open-data dump and rewrites the matching rows into a new CSV that keeps *exactly* the same column order as the official `activity.csv` file.
It lets you run a one-liner such as

```bash
node filter.js 96093
```

to obtain `filtre_96093.csv`, ready for spreadsheet work or further processing.
The raw data come from the monthly “Full” archive published by the FPS Economy ([Économie][1]), while the list of valid sector codes is maintained in the NACE-BEL 2025 table ([Économie][2]). Parsing is done with the battle-tested **csv-parser** streaming package ([npm][3]).

## Quick Start

```bash
# 1. install dependencies
npm install

# 2. download and unzip the latest “Full” archive from the CBE portal
#    and put the CSV files under ./bce_yyyy_mm/
#    (see link in “Data Sources” below)

# 3. generate a filtered CSV
node filter.js 62020      # example: Computerconsultancy

# output ➜ filtre_62020.csv
```

## Usage

```
node filter.js <NACE_CODE>
```

* **`<NACE_CODE>`** – any 5-digit code, e.g. `85600` (educational support activities).
* The script writes a file named `filtre_<NACE_CODE>.csv` in the project root, preserving the header
  `"EntityNumber","ActivityGroup","NaceVersion","NaceCode","Classification"`.

### Changing the data folder

If your CSVs live elsewhere, edit the constant `DIR` at the top of `filter.js`.

## File structure

```
.
├─ bce_2025_05/          # unzipped “Full” archive from CBE
│   ├─ activity.csv
│   └─ ...
├─ filter.js             # CLI script (CommonJS)
├─ package.json
└─ README.md
```

## Data Sources

* **CBE Open Data portal** – monthly “Full” and “Update” ZIP archives. ([Économie][1])
* **NACE-BEL 2025 codes** – official XLSX table. ([Économie][2])
* **Cookbook BCE Open Data** – column descriptions & update policy. ([Économie][4])

## License

Code licensed under MIT.
CBE datasets are free to reuse under the FPS Economy open-data licence; check the cookbook for attribution terms.

[1]: https://economie.fgov.be/en/themes/enterprises/crossroads-bank-enterprises/services-everyone/public-data-available-reuse/cbe-open-data "CBE - Open data | FPS Economy"
[2]: https://economie.fgov.be/en/themes/enterprises/crossroads-bank-enterprises/services-administrations/tables-codes "Tables of codes | FPS Economy"
[3]: https://www.npmjs.com/package/csv-parser "csv-parser - NPM"
[4]: https://economie.fgov.be/sites/default/files/Files/Entreprises/BCE/Cookbook-BCE-Open-Data.pdf "[PDF] Cookbook BCE Open Data Version R015.00 - FOD Economie"
