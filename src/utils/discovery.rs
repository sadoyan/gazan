use crate::utils::consul;
use crate::utils::filewatch;
use crate::utils::parceyaml::Configuration;
use crate::web::webserver;
use async_trait::async_trait;
use futures::channel::mpsc::Sender;

pub struct FromFileProvider {
    pub path: String,
}
pub struct APIUpstreamProvider {
    pub address: String,
}

pub struct ConsulProvider {
    pub path: String,
}

#[async_trait]
pub trait Discovery {
    async fn start(&self, tx: Sender<Configuration>);
}

#[async_trait]
impl Discovery for APIUpstreamProvider {
    async fn start(&self, toreturn: Sender<Configuration>) {
        webserver::run_server(self.address.clone(), toreturn).await;
    }
}

#[async_trait]
impl Discovery for FromFileProvider {
    async fn start(&self, tx: Sender<Configuration>) {
        tokio::spawn(filewatch::start(self.path.clone(), tx.clone()));
    }
}

#[async_trait]
impl Discovery for ConsulProvider {
    async fn start(&self, tx: Sender<Configuration>) {
        tokio::spawn(consul::start(self.path.clone(), tx.clone()));
    }
}
