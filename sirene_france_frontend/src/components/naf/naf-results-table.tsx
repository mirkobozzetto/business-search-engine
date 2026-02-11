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
import { Badge } from "@/components/ui/badge";
import type { NafCode } from "@/types/naf";

interface NafResultsTableProps {
  codes: NafCode[];
}

export function NafResultsTable({ codes }: NafResultsTableProps) {
  if (codes.length === 0) return null;

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[120px]">Code</TableHead>
          <TableHead>Libellé</TableHead>
          <TableHead className="w-[100px]">Section</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {codes.map((naf) => (
          <TableRow key={naf.code}>
            <TableCell>
              <Link
                href={`/naf/${naf.code}`}
                className="font-mono font-medium text-primary hover:underline"
              >
                {naf.code}
              </Link>
            </TableCell>
            <TableCell>{naf.label}</TableCell>
            <TableCell>
              <Link href={`/naf/section/${naf.section_code}`}>
                <Badge variant="secondary">{naf.section_code}</Badge>
              </Link>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
