package lula.labels

import rego.v1

has_lula_label(pod) if {
    pod.metadata.labels.lula == "true"
}
