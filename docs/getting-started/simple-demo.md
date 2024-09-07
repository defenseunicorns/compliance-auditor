# Simple Demo

The following simple demo will step through a process to validate and evaluate Kubernetes cluster resources for a simple `component-definition`. The following pre-requisites are required to successfully run this demo:

### Pre-Requisites

* Lula installed
* Kubectl
* A running Kubernetes cluster
    - Kind
        - `kind create cluster -n lula-test`
    - K3d
        - `k3d cluster create lula-test`

### Steps

1. Clone the Lula repository to your local machine and change into the `lula/demo/simple` directory

    ```shell
    git clone https://github.com/defenseunicorns/lula.git && cd lula/demo/simple
    ```

2. Apply the `namespace.yaml` file to create a namespace for the demo

    ```shell
    kubectl apply -f namespace.yaml
    ```

3. Apply the `pod.fail.yaml` to create a pod in your cluster

    ```shell
    kubectl apply -f pod.fail.yaml
    ```

4. The `oscal-component-opa.yaml` is a simple OSCAL Component Definition model which establishes a sample control <> validation mapping. The validation that provides a `satisfied` or `not-satisfied` result to the control simply checks if the required label value is set for `pod.fail`. Run the following command to `validate` the component given by the failing pod:

    ```shell
    lula validate -f oscal-component-opa.yaml
    ```

    The output in your terminal should inform you that the single control validated is `not-satisfied`:

    ```shell
    NOTE  Saving log file to
        /var/folders/6t/7mh42zsx6yv_3qzw2sfyh5f80000gn/T/lula-2024-07-08-10-22-57-840485568.log
    •  changing cwd to .
    
    🔍 Collecting Requirements and Validations   
    •  Found 1 Implemented Requirements
    •  Found 1 runnable Lula Validations
    
    📐 Running Validations   
    ✔  Running validation a7377430-2328-4dc4-a9e2-b3f31dc1dff9 -> evaluated -> not-satisfied                                            
    
    💡 Findings   
    •  UUID: c80f76c5-4c86-4773-91a3-ece127f3d55a
    •  Status: not-satisfied
    •  OSCAL artifact written to: assessment-results.yaml
    ```

    This will also produce an `assessment-results` model - review the findings and observations:

    ```yaml
      assessment-results:
        results:
            - description: Assessment results for performing Validations with Lula version v0.4.1-1-gc270673
            findings:
                - description: This control validates that the demo-pod pod in the validation-test namespace contains the required pod label foo=bar in order to establish compliance.
                related-observations:
                    - observation-uuid: f03ffcd9-c18d-40bf-85f5-d0b1a8195ddb
                target:
                    status:
                    state: not-satisfied
                    target-id: ID-1
                    type: objective-id
                title: 'Validation Result - Component:A9D5204C-7E5B-4C43-BD49-34DF759B9F04 / Control Implementation: A584FEDC-8CEA-4B0C-9F07-85C2C4AE751A / Control:  ID-1'
                uuid: c80f76c5-4c86-4773-91a3-ece127f3d55a
            observations:
                - collected: 2024-07-08T10:22:57.219213-04:00
                description: |
                    [TEST]: a7377430-2328-4dc4-a9e2-b3f31dc1dff9 - lula-validation
                methods:
                    - TEST
                relevant-evidence:
                    - description: |
                        Result: not-satisfied
                uuid: f03ffcd9-c18d-40bf-85f5-d0b1a8195ddb
            props:
                - name: threshold
                ns: https://docs.lula.dev/oscal/ns
                value: "false"
            reviewed-controls:
                control-selections:
                - description: Controls Assessed by Lula
                    include-controls:
                    - control-id: ID-1
                description: Controls validated
                remarks: Validation performed may indicate full or partial satisfaction
            start: 2024-07-08T10:22:57.219371-04:00
            title: Lula Validation Result
            uuid: f9ae56df-8709-49be-a230-2d3962bbd5f9
        uuid: 5bf89b23-6172-47c9-9d1c-d308fa543d61
    ```

5. Now, apply the `pod.pass.yaml` file to your cluster to configure the pod to pass compliance validation:

    ```shell
    kubectl apply -f pod.pass.yaml
    ```

6. Run the following command in the `lula` directory:

    ```shell
    lula validate -f oscal-component-opa.yaml
    ```

    The output should now show the pod as passing the compliance requirement:

    ```shell
    NOTE  Saving log file to
        /var/folders/6t/7mh42zsx6yv_3qzw2sfyh5f80000gn/T/lula-2024-07-08-10-25-47-3097295143.log
    •  changing cwd to .
    
    🔍 Collecting Requirements and Validations   
    •  Found 1 Implemented Requirements
    •  Found 1 runnable Lula Validations
    
    📐 Running Validations   
    ✔  Running validation a7377430-2328-4dc4-a9e2-b3f31dc1dff9 -> evaluated -> satisfied                                                
    
    💡 Findings   
    •  UUID: 5a991d1f-745e-4acb-9435-373174816fcc
    •  Status: satisfied
    •  OSCAL artifact written to: assessment-results.yaml
    ```

    This will append to the assessment-results file with a new result:

    ```yaml
      - description: Assessment results for performing Validations with Lula version v0.4.1-1-gc270673
      findings:
        - description: This control validates that the demo-pod pod in the validation-test namespace contains the required pod label foo=bar in order to establish compliance.
          related-observations:
            - observation-uuid: a1d55b82-c63f-47da-8fab-87ae801357ac
          target:
            status:
              state: satisfied
            target-id: ID-1
            type: objective-id
          title: 'Validation Result - Component:A9D5204C-7E5B-4C43-BD49-34DF759B9F04 / Control Implementation: A584FEDC-8CEA-4B0C-9F07-85C2C4AE751A / Control:  ID-1'
          uuid: 5a991d1f-745e-4acb-9435-373174816fcc
      observations:
        - collected: 2024-07-08T10:25:47.633634-04:00
          description: |
            [TEST]: a7377430-2328-4dc4-a9e2-b3f31dc1dff9 - lula-validation
          methods:
            - TEST
          relevant-evidence:
            - description: |
                Result: satisfied
          uuid: a1d55b82-c63f-47da-8fab-87ae801357ac
      props:
        - name: threshold
          ns: https://docs.lula.dev/oscal/ns
          value: "false"
      reviewed-controls:
        control-selections:
          - description: Controls Assessed by Lula
            include-controls:
              - control-id: ID-1
        description: Controls validated
        remarks: Validation performed may indicate full or partial satisfaction
      start: 2024-07-08T10:25:47.6341-04:00
      title: Lula Validation Result
      uuid: a9736e32-700d-472f-96a3-4dacf36fa9ce
    ```

7. Now that two assessment-results are established, the `threshold` can be evaluated. Perform an `evaluate` to compare the old and new state of the cluster:
    ```shell
    lula evaluate -f assessment-results.yaml
    ```

    The output will show that now the new threshold for the system assessment is the more _compliant_ evaluation of the control - i.e., the `satisfied` value of the Control ID-1 is the threshold.
    ```shell
     NOTE  Saving log file to
        /var/folders/6t/7mh42zsx6yv_3qzw2sfyh5f80000gn/T/lula-2024-07-08-10-29-53-4238890270.log
    •  New passing finding Target-Ids:                                                                                                                                                                                                                                                                          
    •  ID-1                                                                                                                                                                                                                                                                                                     
    •  New threshold identified - threshold will be updated to result a9736e32-700d-472f-96a3-4dacf36fa9ce                                                                                                                                                                                                      
    •  Evaluation Passed Successfully                                                                                                                                                                                                                                                                           
    ✔  Evaluating Assessment Results f9ae56df-8709-49be-a230-2d3962bbd5f9 against a9736e32-700d-472f-96a3-4dacf36fa9ce                                                                                                                                                                                          
    •  OSCAL artifact written to: ./assessment-results.yaml
    ```