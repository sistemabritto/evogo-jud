# evogo-jud (fork do Evolution Go: Sistema Britto / Cramer)

`main` é a branch única deste fork: Evolution Go 0.7.2 (9337afc) mais os ajustes abaixo. `deploy/sistemabritto` acompanha `main`.

- Prévia de vídeo (miniatura, duração e dimensões) no envio de vídeo.
- Versão anunciada pelo cliente WhatsApp 1052556849 na whatsmeow fixada (erro 405 "Client outdated"); aplicada no Dockerfile. Pode ficar obsoleta se o WhatsApp rejeitar a versão de novo.
- Pedido de histórico sob demanda enviado ao próprio aparelho e captura opcional (HISTORY_SYNC_DUMP_DIR / HISTORY_SYNC_DUMP_INSTANCES).
- Um único `sqlstore.Container` por banco, com pool limitado. Antes, cada reconexão de instância abria um pool novo que nunca era fechado, e as conexões ociosas enchiam o Postgres até `max_connections`.

A Action publica `excarplex/evogo-jud:<sha>` e `:latest` (amd64) a cada push em `main`; por compatibilidade, também `excarplex/evolution-go:sistemabritto` e `:sistemabritto-<sha>` (as mesmas tags saem nos dois nomes). Secrets: `DOCKERHUB_USERNAME` e `DOCKERHUB_TOKEN`. Teste a imagem antes do rollout; para produção, use a tag com SHA.
