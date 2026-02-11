"use client";

import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

interface CompanySearchFormProps {
  denomination: string;
  nafCode: string;
  codePostal: string;
  commune: string;
  etatAdministratif: string;
  onDenominationChange: (v: string) => void;
  onNafCodeChange: (v: string) => void;
  onCodePostalChange: (v: string) => void;
  onCommuneChange: (v: string) => void;
  onEtatAdministratifChange: (v: string) => void;
  onSearch: () => void;
  onReset: () => void;
}

export function CompanySearchForm({
  denomination,
  nafCode,
  codePostal,
  commune,
  etatAdministratif,
  onDenominationChange,
  onNafCodeChange,
  onCodePostalChange,
  onCommuneChange,
  onEtatAdministratifChange,
  onSearch,
  onReset,
}: CompanySearchFormProps) {
  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <div className="space-y-2">
          <label className="text-sm font-medium">Dénomination</label>
          <Input
            value={denomination}
            onChange={(e) => onDenominationChange(e.target.value)}
            placeholder="Nom de l'entreprise"
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">Code NAF</label>
          <Input
            value={nafCode}
            onChange={(e) => onNafCodeChange(e.target.value)}
            placeholder="Ex: 62.01Z"
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">Code postal</label>
          <Input
            value={codePostal}
            onChange={(e) => onCodePostalChange(e.target.value)}
            placeholder="Ex: 75001"
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">Commune</label>
          <Input
            value={commune}
            onChange={(e) => onCommuneChange(e.target.value)}
            placeholder="Ex: PARIS"
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">Statut</label>
          <Select value={etatAdministratif} onValueChange={onEtatAdministratifChange}>
            <SelectTrigger>
              <SelectValue placeholder="Tous" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Tous</SelectItem>
              <SelectItem value="A">Active</SelectItem>
              <SelectItem value="C">Fermée</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>
      <div className="flex gap-2">
        <Button onClick={onSearch}>Rechercher</Button>
        <Button variant="outline" onClick={onReset}>
          Réinitialiser
        </Button>
      </div>
    </div>
  );
}
