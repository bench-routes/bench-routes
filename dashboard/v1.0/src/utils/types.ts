export interface SystemInformation {
  sysinfo: sysinfo;
  node: node;
  os: os;
  kernel: kernel;
  product: product;
  board: board;
  chassis: chassis;
  bios: bios;
  cpu: cpu;
  memory: memory;
  storage: storage[];
  network: network[];
}

interface sysinfo {
  version: string;
  timestamp: string;
}

interface node {
  hostname: string;
  machineid: string;
}

interface os {
  name: string;
  vendor: string;
  version: string;
  architecture: string;
}

interface kernel {
  release: string;
  version: number;
  architecture: string;
}

interface product {
  name: string;
  vendor: string;
  version: string;
}

interface board {
  name: string;
  vendor: string;
  version: string;
}

interface chassis {
  type: number;
  vendor: string;
  version: string;
}

interface bios {
  vendor: string;
  version: string;
  date: string;
}

interface cpu {
  vendor: string;
  model: string;
  cache: number;
  cpus: number;
  cores: number;
  threads: number;
}

interface memory {}

interface storage {
  name: string;
  driver: string;
  vendor: string;
  model: string;
  serial: string;
  size: number;
}

interface network {
  name: string;
  driver: string;
  macaddress: string;
  port: string;
  speed: number;
}

export interface service_states {
  ping: string;
  floodping: string;
  jitter: string;
  monitoring: string;
}

export const HOST_IP = 'http://localhost:9090';
