import { Injectable } from '@nestjs/common';
import type { Secret } from '@prisma/client';
import { SecretRepository } from './secret.repository';
import { AppError } from 'src/common/http/errors/app-error';
import { ErrorDefinitions } from 'src/common/http/errors/error-codes';

@Injectable()
export class SecretService {
  constructor(private readonly repo: SecretRepository) {}

  async create(params: {
    name: string;
    value: string;
    description?: string;
  }): Promise<Secret> {
    return this.repo.create({
      name: params.name,
      value: params.value,
      description: params.description,
    });
  }

  list(): Promise<Secret[]> {
    return this.repo.findMany();
  }

  async get(id: string): Promise<Secret> {
    const secret = await this.repo.findById(id);
    if (!secret) throw AppError.notFound(ErrorDefinitions.SECRET.NOT_FOUND);
    return secret;
  }

  async update(
    id: string,
    patch: {
      name?: string;
      value?: string;
      description?: string;
    },
  ): Promise<Secret> {
    await this.get(id);
    return this.repo.update(id, patch);
  }

  async delete(id: string): Promise<Secret> {
    await this.get(id);
    return this.repo.delete(id);
  }
}
