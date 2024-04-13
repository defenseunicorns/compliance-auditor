# OSCAL Generation

Lula has the potential to provide a codified base for how to generate and maintain OSCAL through automation. This means that with a foundation built - Lula can continue to iterate on the methods for mapping maintaining data the aligns with the intent of OSCAL and the standards/benchmarks involved. 

## Generic Generation Concepts

The generation process for OSCAL artifacts created and maintained by Lula should include the following:
- Specification of fields only maintained by automation (Challenge this one)
- Ability to maintain data that is added through manual interaction

## Component Definition Generation

Current TODO:
- Ability to fetch a catalog from source at some external location (IE CatalogSource flag)
  - Detection of file type (json/yaml)
- Ability to detect an existing output file (IE OutputFile flag)
- Ability to retain data in an existing OutputFile on re-generation
- Ability to detect an OSCAL manifest file (IE InputFile flag)