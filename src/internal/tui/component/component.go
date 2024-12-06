package component

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/defenseunicorns/lula/src/internal/tui/common"
	pkgcommon "github.com/defenseunicorns/lula/src/pkg/common"
	"github.com/defenseunicorns/lula/src/pkg/common/composition"
	"github.com/defenseunicorns/lula/src/pkg/common/oscal"
)

type component struct {
	OscalComponent   *oscalTypes.DefinedComponent
	Uuid, Name, Desc string
	Frameworks       []framework
}

type framework struct {
	OscalFramework *oscalTypes.ControlImplementationSet
	Uuid, Name     string
	Controls       []control
}

type control struct {
	OscalControl *oscalTypes.ImplementedRequirementControlImplementation
	Uuid, Name   string
	Validations  []validationLink
}

func (i control) Title() string       { return i.Name }
func (i control) Description() string { return i.Uuid }
func (i control) FilterValue() string { return i.Name }

type validationLink struct {
	OscalLink  *oscalTypes.Link
	Text       string
	Name       string
	Validation pkgcommon.Validation
}

func (i validationLink) Title() string       { return i.Name }
func (i validationLink) Description() string { return i.Text }
func (i validationLink) FilterValue() string { return i.Name }

func GetComponents(oscalComponent *oscalTypes.ComponentDefinition) []component {
	components := make([]component, 0)

	if oscalComponent != nil {
		// TODO: Add composition, path - assuming the input model is already composed
		resourceStore := composition.NewResourceStoreFromBackMatter(nil, oscalComponent.BackMatter)

		if oscalComponent.Components != nil {
			for cIdx, c := range *oscalComponent.Components {
				// for each component, add the control implementation to the framework
				componentFrameworks := make([]framework, 0)
				if c.ControlImplementations != nil {
					for ctrlImpIdx, controlImpl := range *c.ControlImplementations {
						// get the controls for each framework
						controls := make([]control, 0, len(controlImpl.ImplementedRequirements))
						for reqIdx, implementedRequirement := range controlImpl.ImplementedRequirements {
							// get validations from implementedRequirement.Links
							validationLinks := make([]validationLink, 0)
							if implementedRequirement.Links != nil {
								for _, link := range *implementedRequirement.Links {
									if pkgcommon.IsLulaLink(link) {
										var validation pkgcommon.Validation
										name := "Lula Validation"
										resource, found := resourceStore.Get(pkgcommon.TrimIdPrefix(link.Href)) //pkgcommon.TrimIdPrefix(link.Href)
										if found {
											err := validation.UnmarshalYaml([]byte(resource.Description))
											if err != nil {
												common.PrintToLog("error unmarshalling validation: %v", err)
											}
										}
										if validation.Metadata != nil {
											if validation.Metadata.Name != "" {
												name = validation.Metadata.Name
											}
										}
										// add the lula validation to the validationsLinks array
										validationLinks = append(validationLinks, validationLink{
											OscalLink:  &link, //&(*(*c.ControlImplementations)[ctrlImpIdx].ImplementedRequirements[reqIdx].Links)[linkIdx],
											Text:       link.Text,
											Name:       name,
											Validation: validation,
										})

									}
								}
							}

							controls = append(controls, control{
								OscalControl: &(*c.ControlImplementations)[ctrlImpIdx].ImplementedRequirements[reqIdx], //&implementedRequirement,
								Name:         implementedRequirement.ControlId,
								Uuid:         implementedRequirement.UUID,
								Validations:  validationLinks,
							})
						}

						// sort controls by title
						slices.SortStableFunc(controls, func(a, b control) int {
							return oscal.CompareControlsInt(a.Name, b.Name)
						})

						componentFrameworks = append(componentFrameworks, framework{
							OscalFramework: &(*c.ControlImplementations)[ctrlImpIdx], //&controlImpl,
							Name:           controlImpl.Source,
							Uuid:           controlImpl.UUID,
							Controls:       controls,
						})

						// Add named framework if set
						status, value := oscal.GetProp("framework", oscal.LULA_NAMESPACE, controlImpl.Props)
						if status {
							componentFrameworks = append(componentFrameworks, framework{
								OscalFramework: &(*c.ControlImplementations)[ctrlImpIdx], //&controlImpl,
								Name:           value,
								Uuid:           controlImpl.UUID,
								Controls:       controls,
							})
						}
					}
				}

				// sort componentFrameworks by name
				slices.SortStableFunc(componentFrameworks, func(a, b framework) int {
					return strings.Compare(a.Name, b.Name)
				})

				components = append(components, component{
					OscalComponent: &(*oscalComponent.Components)[cIdx], //&c,
					Uuid:           c.UUID,
					Name:           c.Title,
					Desc:           c.Description,
					Frameworks:     componentFrameworks,
				})
			}

			if len(components) > 0 {
				// sort components by title
				slices.SortStableFunc(components, func(a, b component) int {
					return strings.Compare(a.Name, b.Name)
				})
			}
		}
	}

	return components
}

func getComponentText(component component) string {
	if component.Uuid == "" {
		return "No Components"
	}
	return fmt.Sprintf("%s - %s", component.Name, component.Uuid)
}

func getFrameworkText(framework framework) string {
	if framework.Name == "" {
		return "No Frameworks"
	}
	return framework.Name
}

func getValidationText(validationLink validationLink) string {
	var text strings.Builder
	validation := validationLink.Validation

	important := lipgloss.NewStyle().Bold(true).
		Foreground(common.Special)

	if validation.Metadata != nil {
		text.WriteString(fmt.Sprintf("%s - %s\n", important.Render(validationLink.Name), validation.Metadata.UUID))
	} else {
		text.WriteString(fmt.Sprintf("%s\n", important.Render(validationLink.Name)))
	}
	text.WriteString("\n\n")

	if validation.Domain != nil {
		text.WriteString(fmt.Sprintf("Domain: %s\n", important.Render(validation.Domain.Type)))
		switch validation.Domain.Type {
		case "kubernetes":
			kubeSpec, err := common.ToYamlString(validation.Domain.KubernetesSpec)
			if err != nil {
				common.PrintToLog("error converting kubeSpec to yaml: %v", err)
				kubeSpec = ""
			}
			text.WriteString(kubeSpec)
		case "api":
			apiSpec, err := common.ToYamlString(validation.Domain.ApiSpec)
			if err != nil {
				common.PrintToLog("error converting apiSpec to yaml: %v", err)
				apiSpec = ""
			}
			text.WriteString(apiSpec)
		case "file":
			fileSpec, err := common.ToYamlString(validation.Domain.FileSpec)
			if err != nil {
				common.PrintToLog("error converting fileSpec to yaml: %v", err)
				fileSpec = ""
			}
			text.WriteString(fileSpec)
		}
		text.WriteString("\n\n")
	}

	if validation.Provider != nil {
		text.WriteString(fmt.Sprintf("Provider: %s\n", important.Render(validation.Provider.Type)))
		switch validation.Provider.Type {
		case "opa":
			opaSpec, err := common.ToYamlString(validation.Provider.OpaSpec)
			if err != nil {
				common.PrintToLog("error converting opaSpec to yaml: %v", err)
				opaSpec = ""
			}
			text.WriteString(opaSpec)
		case "kyverno":
			kyvernoSpec, err := common.ToYamlString(validation.Provider.KyvernoSpec)
			if err != nil {
				common.PrintToLog("error converting kyvernoSpec to yaml: %v", err)
				kyvernoSpec = ""
			}
			text.WriteString(kyvernoSpec)
		}
	}

	return text.String()
}
