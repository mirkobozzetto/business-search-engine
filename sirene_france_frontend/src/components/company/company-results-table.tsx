"use client";

import Link from "next/link";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { CompanyStatusBadge } from "./company-status-badge";
import { formatDateFr } from "@/lib/labels";
import type { CompanyResult } from "@/types/company";

interface CompanyResultsTableProps {
  companies: CompanyResult[];
}

export function CompanyResultsTable({ companies }: CompanyResultsTableProps) {
  if (companies.length === 0) return null;

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Dénomination</TableHead>
          <TableHead className="w-[130px]">SIREN</TableHead>
          <TableHead>Code NAF</TableHead>
          <TableHead>Commune</TableHead>
          <TableHead className="w-[100px]">Création</TableHead>
          <TableHead className="w-[80px]">Statut</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {companies.map((company, idx) => (
          <TableRow key={`${company.siren}-${idx}`}>
            <TableCell>
              <Link
                href={`/company/${company.siren}`}
                className="font-medium text-primary hover:underline"
              >
                {company.denomination || company.siren}
              </Link>
              {company.sigle && (
                <span className="ml-2 text-sm text-muted-foreground">
                  ({company.sigle})
                </span>
              )}
            </TableCell>
            <TableCell className="font-mono">{company.siren}</TableCell>
            <TableCell>
              {company.naf_code && (
                <Link
                  href={`/naf/${company.naf_code}`}
                  className="text-sm hover:underline"
                >
                  {company.naf_code}
                </Link>
              )}
            </TableCell>
            <TableCell>
              {company.libelle_commune && (
                <span className="text-sm">
                  {company.code_postal} {company.libelle_commune}
                </span>
              )}
            </TableCell>
            <TableCell className="text-sm">{formatDateFr(company.date_creation)}</TableCell>
            <TableCell>
              <CompanyStatusBadge status={company.etat_administratif} />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
