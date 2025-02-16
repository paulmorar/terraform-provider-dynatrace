resource "dynatrace_log_metrics" "#name#" {
  enabled           = true
  dimensions        = [ "dt.os.type", "dt.entity.process_group" ]
  key               = "log.terraformexample"
  measure           = "ATTRIBUTE"
  measure_attribute = "dt.entity.host"
  query             = "TestMatcherValue"
}