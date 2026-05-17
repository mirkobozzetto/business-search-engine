"use client";

import { useState } from "react";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Input } from "@/components/ui/input";

interface PaginationControlsProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

function buildPageNumbers(current: number, total: number): (number | "ellipsis-start" | "ellipsis-end")[] {
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1);
  }

  const pages: (number | "ellipsis-start" | "ellipsis-end")[] = [1];

  if (current > 3) {
    pages.push("ellipsis-start");
  }

  const start = Math.max(2, current - 1);
  const end = Math.min(total - 1, current + 1);

  for (let i = start; i <= end; i++) {
    pages.push(i);
  }

  if (current < total - 2) {
    pages.push("ellipsis-end");
  }

  pages.push(total);

  return pages;
}

export function PaginationControls({
  currentPage,
  totalPages,
  onPageChange,
}: PaginationControlsProps) {
  const [goToValue, setGoToValue] = useState("");

  if (totalPages <= 1) return null;

  const pages = buildPageNumbers(currentPage, totalPages);

  function handleGoTo() {
    const page = parseInt(goToValue, 10);
    if (page >= 1 && page <= totalPages) {
      onPageChange(page);
      setGoToValue("");
    }
  }

  return (
    <div className="flex flex-col items-center gap-3 sm:flex-row sm:justify-between">
      <Pagination>
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious
              onClick={() => onPageChange(Math.max(1, currentPage - 1))}
              className={currentPage <= 1 ? "pointer-events-none opacity-50" : "cursor-pointer"}
            />
          </PaginationItem>
          {pages.map((page) =>
            typeof page === "string" ? (
              <PaginationItem key={page}>
                <PaginationEllipsis />
              </PaginationItem>
            ) : (
              <PaginationItem key={page}>
                <PaginationLink
                  onClick={() => onPageChange(page)}
                  isActive={page === currentPage}
                  className="cursor-pointer"
                >
                  {page}
                </PaginationLink>
              </PaginationItem>
            )
          )}
          <PaginationItem>
            <PaginationNext
              onClick={() => onPageChange(Math.min(totalPages, currentPage + 1))}
              className={currentPage >= totalPages ? "pointer-events-none opacity-50" : "cursor-pointer"}
            />
          </PaginationItem>
        </PaginationContent>
      </Pagination>

      <div className="flex items-center gap-2 text-sm">
        <span className="text-muted-foreground whitespace-nowrap">Aller à</span>
        <Input
          type="number"
          min={1}
          max={totalPages}
          value={goToValue}
          onChange={(e) => setGoToValue(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleGoTo()}
          className="h-8 w-16 text-center"
          placeholder={String(currentPage)}
        />
        <span className="text-muted-foreground whitespace-nowrap">/ {totalPages}</span>
      </div>
    </div>
  );
}
