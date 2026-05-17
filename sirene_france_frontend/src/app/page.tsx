import { SearchBarHero } from "@/components/layout/search-bar-hero";

export default function HomePage() {
  return (
    <div className="flex flex-col items-center justify-center gap-8 py-24">
      <div className="text-center space-y-4">
        <h1 className="text-4xl font-bold tracking-tight">
          SIRENE France
        </h1>
        <p className="text-lg text-muted-foreground max-w-xl">
          Explorez la base SIRENE : 36 millions d'établissements,
          24 millions d'unités légales, 732 codes NAF.
        </p>
      </div>
      <SearchBarHero />
      <div className="flex gap-8 text-center text-sm text-muted-foreground">
        <div>
          <p className="text-2xl font-bold text-foreground">36M</p>
          <p>Établissements</p>
        </div>
        <div>
          <p className="text-2xl font-bold text-foreground">24M</p>
          <p>Unités légales</p>
        </div>
        <div>
          <p className="text-2xl font-bold text-foreground">732</p>
          <p>Codes NAF</p>
        </div>
      </div>
    </div>
  );
}
