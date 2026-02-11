import Link from "next/link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { CompanyStatusBadge } from "./company-status-badge";
import { getCategorieJuridiqueLabel, getTrancheEffectifsLabel, formatDateFr } from "@/lib/labels";
import type { CompanyResult } from "@/types/company";

interface CompanyDetailProps {
  company: CompanyResult;
}

function DetailRow({ label, value }: { label: string; value?: string | null }) {
  if (!value) return null;
  return (
    <div className="grid grid-cols-3 gap-4 py-2">
      <span className="text-sm font-medium text-muted-foreground">{label}</span>
      <span className="col-span-2 text-sm">{value}</span>
    </div>
  );
}

export function CompanyDetail({ company }: CompanyDetailProps) {
  const address = [
    company.numero_voie,
    company.type_voie,
    company.libelle_voie,
  ]
    .filter(Boolean)
    .join(" ");

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="text-xl">
              {company.denomination || company.siren}
            </CardTitle>
            <CompanyStatusBadge status={company.etat_administratif} />
          </div>
          {company.sigle && (
            <p className="text-sm text-muted-foreground">{company.sigle}</p>
          )}
        </CardHeader>
        <CardContent className="space-y-1">
          <DetailRow label="SIREN" value={company.siren} />
          <DetailRow label="SIRET" value={company.siret} />
          <DetailRow label="Date de création" value={formatDateFr(company.date_creation)} />
          <DetailRow label="Catégorie juridique" value={company.categorie_juridique ? `${getCategorieJuridiqueLabel(company.categorie_juridique)} (${company.categorie_juridique})` : undefined} />
          <DetailRow label="Catégorie entreprise" value={company.categorie_entreprise} />
          <DetailRow label="Tranche effectifs" value={getTrancheEffectifsLabel(company.tranche_effectifs)} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Activité</CardTitle>
        </CardHeader>
        <CardContent className="space-y-1">
          {company.naf_code && (
            <div className="grid grid-cols-3 gap-4 py-2">
              <span className="text-sm font-medium text-muted-foreground">Code NAF</span>
              <span className="col-span-2 text-sm">
                <Link
                  href={`/naf/${company.naf_code}`}
                  className="text-primary hover:underline"
                >
                  {company.naf_code}
                </Link>
                {company.naf_label && ` - ${company.naf_label}`}
              </span>
            </div>
          )}
          {company.enseigne && (
            <DetailRow label="Enseigne" value={company.enseigne} />
          )}
        </CardContent>
      </Card>

      {(address || company.code_postal || company.libelle_commune) && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Adresse</CardTitle>
          </CardHeader>
          <CardContent className="space-y-1">
            <DetailRow label="Voie" value={address || undefined} />
            <DetailRow label="Code postal" value={company.code_postal} />
            <DetailRow label="Commune" value={company.libelle_commune} />
          </CardContent>
        </Card>
      )}

      {(company.email || company.telephone || company.website) && (
        <>
          <Separator />
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Contact</CardTitle>
            </CardHeader>
            <CardContent className="space-y-1">
              <DetailRow label="Email" value={company.email} />
              <DetailRow label="Téléphone" value={company.telephone} />
              <DetailRow label="Site web" value={company.website} />
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
}
