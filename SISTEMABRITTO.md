# Imagem Sistema Britto

A branch `deploy/sistemabritto` parte do Evolution Go 0.7.2 (9337afc). Inclui o patch de prévia de vídeo desenvolvido para Cramer e atualiza a versão anunciada pelo cliente WhatsApp de 1042386815 para 1052556849 na dependência whatsmeow fixada. Este ajuste pode se tornar obsoleto se o WhatsApp rejeitar novamente a versão.

A Action publica `excarplex/evolution-go:sistemabritto` e uma tag com SHA. Configure os secrets `DOCKERHUB_USERNAME` e `DOCKERHUB_TOKEN` (token com permissão de escrita) no repositório. A imagem deve ser testada antes do rollout. Reconexão não garante importação do histórico do WhatsApp.
