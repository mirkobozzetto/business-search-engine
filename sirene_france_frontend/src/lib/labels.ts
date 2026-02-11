export const CATEGORIES_JURIDIQUES: Record<string, string> = {
  "1000": "Entrepreneur individuel",
  "5498": "EURL",
  "5499": "SASU",
  "5505": "SA à conseil d'administration",
  "5510": "SA à directoire",
  "5710": "SAS",
  "5720": "SARL",
  "5785": "SELAS",
  "6540": "SCI",
  "9220": "Association déclarée",
  "9221": "Association déclarée d'insertion",
  "9300": "Fondation",
};

export const TRANCHES_EFFECTIFS: Record<string, string> = {
  "NN": "Non renseignée",
  "00": "0 salarié",
  "01": "1 à 2 salariés",
  "02": "3 à 5 salariés",
  "03": "6 à 9 salariés",
  "11": "10 à 19 salariés",
  "12": "20 à 49 salariés",
  "21": "50 à 99 salariés",
  "22": "100 à 199 salariés",
  "31": "200 à 249 salariés",
  "32": "250 à 499 salariés",
  "41": "500 à 999 salariés",
  "42": "1 000 à 1 999 salariés",
  "51": "2 000 à 4 999 salariés",
  "52": "5 000 à 9 999 salariés",
  "53": "10 000 salariés et plus",
};

export function getCategorieJuridiqueLabel(code?: string): string {
  if (!code) return "";
  return CATEGORIES_JURIDIQUES[code] || code;
}

export function getTrancheEffectifsLabel(code?: string): string {
  if (!code) return "";
  return TRANCHES_EFFECTIFS[code] || code;
}

export function formatDateFr(dateStr?: string): string {
  if (!dateStr) return "";
  const parts = dateStr.split("-");
  if (parts.length !== 3) return dateStr;
  return `${parts[2]}/${parts[1]}/${parts[0]}`;
}
