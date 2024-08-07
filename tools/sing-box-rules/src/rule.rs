use std::str::FromStr;

use ipnet::IpNet;
use serde::de::{Error, Unexpected};
use serde::{Deserialize, Serialize};
use serde_yaml::Error as YamlError;
use serde_yaml::Result as YamlResult;

#[derive(Serialize, Deserialize)]
pub enum ClashRuleType {
    Domain,
    DomainSuffix,
    DomainKeyword,
    IPCidr,
    IPv6Cidr,
    SrcIPCidr,
    SrcPort,
    DstPort,
    ProcessName,
    ProcessPath,
}

impl FromStr for ClashRuleType {
    type Err = YamlError;

    fn from_str(s: &str) -> YamlResult<Self> {
        let value = match s {
            "DOMAIN" => ClashRuleType::Domain,
            "DOMAIN-SUFFIX" => ClashRuleType::DomainSuffix,
            "DOMAIN-KEYWORD" => ClashRuleType::DomainKeyword,
            "IP-CIDR" => ClashRuleType::IPCidr,
            "IP-CIDR6" => ClashRuleType::IPv6Cidr,
            "SRC-IP-CIDR" => ClashRuleType::SrcIPCidr,
            "SRC-PORT" => ClashRuleType::SrcPort,
            "DST-PORT" => ClashRuleType::DstPort,
            "PROCESS-NAME" => ClashRuleType::ProcessName,
            "PROCESS-PATH" => ClashRuleType::ProcessPath,
            _ => {
                return Err(YamlError::unknown_field(
                    s,
                    &[
                        "DOMAIN",
                        "DOMAIN-SUFFIX",
                        "DOMAIN-KEYWORD",
                        "IP-CIDR",
                        "IP-CIDR6",
                        "SRC-IP-CIDR",
                        "SRC-PORT",
                        "DST-PORT",
                        "PROCESS-NAME",
                        "PROCESS-PATH",
                    ],
                ))
            }
        };

        Ok(value)
    }
}

#[derive(Deserialize)]
pub struct ClashRule {
    pub rule_type: ClashRuleType,
    pub value: String,
}

impl ClashRule {
    pub fn parse(s: &str) -> YamlResult<Self> {
        let values: Vec<&str> = s.split(",").collect();
        let value = if values.len() >= 2 {
            Self {
                rule_type: ClashRuleType::from_str(values[0])?,
                value: values[1].to_string(),
            }
        } else if values.len() == 1 {
            if let Ok(ip) = IpNet::from_str(values[0]) {
                Self {
                    rule_type: ClashRuleType::IPCidr,
                    value: ip.to_string(),
                }
            } else {
                return Err(YamlError::invalid_type(
                    Unexpected::Str(values[0]),
                    &"IP-CIDR",
                ));
            }
        } else {
            return Err(YamlError::invalid_length(0, &"length >= 1"));
        };

        Ok(value)
    }

    pub fn value_typed<T: FromStr>(&self) -> Option<T> {
        T::from_str(&self.value).ok()
    }
}

#[derive(Deserialize)]
pub struct ClashRuleSet {
    pub payload: Vec<String>,
}

impl ClashRuleSet {
    pub fn new() -> Self {
        Self {
            payload: Vec::new(),
        }
    }

    pub fn from_str(s: &str) -> YamlResult<Self> {
        serde_yaml::from_str(s)
    }

    fn payload(&self) -> Vec<ClashRule> {
        self.payload
            .iter()
            .flat_map(|s| ClashRule::parse(s).ok())
            .collect()
    }
}

#[derive(Serialize)]
pub struct BoxRule {
    #[serde(skip_serializing_if = "Vec::is_empty")]
    pub domain: Vec<String>,
    #[serde(skip_serializing_if = "Vec::is_empty")]
    pub domain_suffix: Vec<String>,
    #[serde(skip_serializing_if = "Vec::is_empty")]
    pub domain_keyword: Vec<String>,
    #[serde(skip_serializing_if = "Vec::is_empty")]
    pub source_ip_cidr: Vec<String>,
    #[serde(skip_serializing_if = "Vec::is_empty")]
    pub ip_cidr: Vec<String>,
    #[serde(skip_serializing_if = "Vec::is_empty")]
    pub source_port: Vec<i32>,
    #[serde(skip_serializing_if = "Vec::is_empty")]
    pub port: Vec<i32>,
    #[serde(skip_serializing_if = "Vec::is_empty")]
    pub process_name: Vec<String>,
    #[serde(skip_serializing_if = "Vec::is_empty")]
    pub process_path: Vec<String>,
}

impl BoxRule {
    pub fn new() -> Self {
        Self {
            domain: Vec::new(),
            domain_suffix: Vec::new(),
            domain_keyword: Vec::new(),
            source_ip_cidr: Vec::new(),
            ip_cidr: Vec::new(),
            source_port: Vec::new(),
            port: Vec::new(),
            process_name: Vec::new(),
            process_path: Vec::new(),
        }
    }

    pub fn from_clash(value: &ClashRuleSet) -> Self {
        let mut rules = BoxRule::new();
        for rule in value.payload() {
            match rule.rule_type {
                ClashRuleType::Domain => rules.domain.push(rule.value),
                ClashRuleType::DomainSuffix => rules.domain_suffix.push(rule.value),
                ClashRuleType::DomainKeyword => rules.domain_keyword.push(rule.value),
                ClashRuleType::IPCidr => rules.ip_cidr.push(rule.value),
                ClashRuleType::IPv6Cidr => rules.ip_cidr.push(rule.value),
                ClashRuleType::SrcIPCidr => rules.source_ip_cidr.push(rule.value),
                ClashRuleType::SrcPort => match rule.value_typed() {
                    Some(v) => rules.source_port.push(v),
                    None => eprintln!("Unknown SrcPort: {}", rule.value),
                },
                ClashRuleType::DstPort => match rule.value_typed() {
                    Some(v) => rules.source_port.push(v),
                    None => eprintln!("Unknown DstPort: {}", rule.value),
                },
                ClashRuleType::ProcessName => rules.process_name.push(rule.value),
                ClashRuleType::ProcessPath => rules.process_path.push(rule.value),
            }
        }

        rules
    }
}

#[derive(Serialize)]
pub struct BoxRuleSet {
    pub version: i32,
    pub rules: Vec<BoxRule>,
}

impl BoxRuleSet {
    pub fn new(rule: BoxRule) -> Self {
        Self {
            version: 1,
            rules: vec![rule],
        }
    }

    pub fn from_clash(value: &ClashRuleSet) -> Self {
        Self::new(BoxRule::from_clash(value))
    }

    pub fn to_string_pretty(&self) -> String {
        serde_json::to_string_pretty(self).unwrap()
    }
}
