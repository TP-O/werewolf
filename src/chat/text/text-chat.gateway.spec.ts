import { Test, TestingModule } from '@nestjs/testing';
import { TextChatGateway } from './text-chat.gateway';

describe('TextChatGateway', () => {
  let gateway: TextChatGateway;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [TextChatGateway],
    }).compile();

    gateway = module.get<TextChatGateway>(TextChatGateway);
  });

  it('should be defined', () => {
    expect(gateway).toBeDefined();
  });
});
