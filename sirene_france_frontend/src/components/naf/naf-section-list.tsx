"use client";

import Link from "next/link";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Badge } from "@/components/ui/badge";
import type { NafSection } from "@/types/naf";

interface NafSectionListProps {
  sections: NafSection[];
}

export function NafSectionList({ sections }: NafSectionListProps) {
  if (sections.length === 0) return null;

  return (
    <Accordion type="single" collapsible className="w-full">
      {sections.map((section) => (
        <AccordionItem key={section.code} value={section.code}>
          <AccordionTrigger className="text-left">
            <div className="flex items-center gap-3">
              <Badge variant="outline" className="font-mono">
                {section.code}
              </Badge>
              <span>{section.label}</span>
              <span className="text-sm text-muted-foreground">
                ({section.count} codes)
              </span>
            </div>
          </AccordionTrigger>
          <AccordionContent>
            <Link
              href={`/naf/section/${section.code}`}
              className="inline-block text-sm text-primary hover:underline"
            >
              Voir les {section.count} codes de la section {section.code}
            </Link>
          </AccordionContent>
        </AccordionItem>
      ))}
    </Accordion>
  );
}
