"use client";

import { use } from "react";
import Link from "next/link";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { CompanyDetail } from "@/components/company/company-detail";
import { LoadingSkeleton } from "@/components/shared/loading-skeleton";
import { EmptyState } from "@/components/shared/empty-state";
import { useCompanyLookup } from "@/hooks/use-company-lookup";

export default function CompanyDetailPage({ params }: { params: Promise<{ identifier: string }> }) {
  const { identifier } = use(params);
  const { data, isLoading, error } = useCompanyLookup(identifier);

  const company = data?.data;

  return (
    <div className="space-y-6">
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink asChild>
              <Link href="/search">Recherche</Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbPage>{identifier}</BreadcrumbPage>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>

      {isLoading ? (
        <LoadingSkeleton rows={8} />
      ) : error ? (
        <EmptyState
          title="Erreur lors du chargement"
          description={error.message}
        />
      ) : company ? (
        <CompanyDetail company={company} />
      ) : (
        <EmptyState
          title="Entreprise introuvable"
          description={`Aucune entreprise trouvée pour l'identifiant ${identifier}`}
        />
      )}
    </div>
  );
}
