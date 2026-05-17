"use client";

import { use } from "react";
import Link from "next/link";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { NafResultsTable } from "@/components/naf/naf-results-table";
import { LoadingSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { useNafSectionCodes } from "@/hooks/use-naf-section-codes";

export default function NafSectionPage({ params }: { params: Promise<{ code: string }> }) {
  const { code } = use(params);
  const { data, isLoading } = useNafSectionCodes(code);

  const codes = data?.data || [];
  const sectionLabel = codes.length > 0 ? codes[0].section_label : "";

  return (
    <div className="space-y-6">
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href="/naf">Codes NAF</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>Section {code}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      <Card>
        <CardHeader>
          <CardTitle>
            Section {code}
            {sectionLabel && ` - ${sectionLabel}`}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <LoadingSkeleton rows={10} />
          ) : codes.length > 0 ? (
            <>
              <p className="mb-4 text-sm text-muted-foreground">
                {codes.length} codes NAF dans cette section
              </p>
              <NafResultsTable codes={codes} />
            </>
          ) : (
            <EmptyState title="Aucun code NAF dans cette section" />
          )}
        </CardContent>
      </Card>
    </div>
  );
}
