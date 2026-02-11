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
  siren: string;
  siret: string;
  denomination: string;
  nafCode: string;
  codePostal: string;
  commune: string;
  etatAdministratif: string;
  dateCreationFrom: string;
  dateCreationTo: string;
  categorieJuridique: string;
  trancheEffectifs: string;
  onSirenChange: (v: string) => void;
  onSiretChange: (v: string) => void;
  onDenominationChange: (v: string) => void;
  onNafCodeChange: (v: string) => void;
  onCodePostalChange: (v: string) => void;
  onCommuneChange: (v: string) => void;
  onEtatAdministratifChange: (v: string) => void;
  onDateCreationFromChange: (v: string) => void;
  onDateCreationToChange: (v: string) => void;
  onCategorieJuridiqueChange: (v: string) => void;
  onTrancheEffectifsChange: (v: string) => void;
  onSearch: () => void;
  onReset: () => void;
}

export function CompanySearchForm({
  siren,
  siret,
  denomination,
  nafCode,
  codePostal,
  commune,
  etatAdministratif,
  dateCreationFrom,
  dateCreationTo,
  categorieJuridique,
  trancheEffectifs,
  onSirenChange,
  onSiretChange,
  onDenominationChange,
  onNafCodeChange,
  onCodePostalChange,
  onCommuneChange,
  onEtatAdministratifChange,
  onDateCreationFromChange,
  onDateCreationToChange,
  onCategorieJuridiqueChange,
  onTrancheEffectifsChange,
  onSearch,
  onReset,
}: CompanySearchFormProps) {
  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
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
        <div className="space-y-2">
          <label className="text-sm font-medium">Date de création (depuis)</label>
          <Input
            type="date"
            value={dateCreationFrom}
            onChange={(e) => onDateCreationFromChange(e.target.value)}
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">Date de création (jusqu'à)</label>
          <Input
            type="date"
            value={dateCreationTo}
            onChange={(e) => onDateCreationToChange(e.target.value)}
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">Catégorie juridique</label>
          <Select value={categorieJuridique} onValueChange={onCategorieJuridiqueChange}>
            <SelectTrigger>
              <SelectValue placeholder="Toutes" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Toutes</SelectItem>
              <SelectItem value="1000">Entrepreneur individuel</SelectItem>
              <SelectItem value="5498">EURL</SelectItem>
              <SelectItem value="5499">SASU</SelectItem>
              <SelectItem value="5505">SA</SelectItem>
              <SelectItem value="5710">SAS</SelectItem>
              <SelectItem value="5720">SARL</SelectItem>
              <SelectItem value="6540">SCI</SelectItem>
              <SelectItem value="9220">Association déclarée</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">Tranche effectifs</label>
          <Select value={trancheEffectifs} onValueChange={onTrancheEffectifsChange}>
            <SelectTrigger>
              <SelectValue placeholder="Toutes" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Toutes</SelectItem>
              <SelectItem value="00">0 salarié</SelectItem>
              <SelectItem value="01">1-2</SelectItem>
              <SelectItem value="02">3-5</SelectItem>
              <SelectItem value="03">6-9</SelectItem>
              <SelectItem value="11">10-19</SelectItem>
              <SelectItem value="12">20-49</SelectItem>
              <SelectItem value="21">50-99</SelectItem>
              <SelectItem value="22">100-199</SelectItem>
              <SelectItem value="31">200-249</SelectItem>
              <SelectItem value="32">250-499</SelectItem>
              <SelectItem value="41">500-999</SelectItem>
              <SelectItem value="42">1000-1999</SelectItem>
              <SelectItem value="51">2000-4999</SelectItem>
              <SelectItem value="52">5000-9999</SelectItem>
              <SelectItem value="53">10000+</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">SIREN</label>
          <Input
            value={siren}
            onChange={(e) => onSirenChange(e.target.value)}
            placeholder="9 chiffres"
            maxLength={9}
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">SIRET</label>
          <Input
            value={siret}
            onChange={(e) => onSiretChange(e.target.value)}
            placeholder="14 chiffres"
            maxLength={14}
          />
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
